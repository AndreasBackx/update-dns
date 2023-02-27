use std::fs;
use std::net::Ipv4Addr;
use std::path::PathBuf;

use anyhow::{bail, Context, Result};
use clap::Parser;
use clap::{arg, command};
use cloudflare::endpoints::{dns, zone};
use cloudflare::framework::async_api::Client;
use cloudflare::framework::auth::Credentials;
use cloudflare::framework::{Environment, HttpApiClientConfig};
use tracing::{debug, info};

fn default_last_ip_path() -> String {
    let xdg_cache_home = std::env::var("XDG_CACHE_HOME").unwrap_or_else(|_| {
        let home = std::env::var("HOME").unwrap();
        format!("{home}/.cache")
    });
    format!("{xdg_cache_home}/update-dns/last_ip")
}

/// Simple program to greet a person
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Entire domain to update, e.g.: "subdomain.example.com".
    #[arg(short, long)]
    domain: String,

    /// Where to store the last IP address, if parent directories do not exist,
    /// they will be created.
    #[arg(short, long, default_value_t = default_last_ip_path())]
    last_ip_path: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();
    let args = Args::parse();

    let path: PathBuf = args.last_ip_path.into();
    let path_parent = path
        .parent()
        .with_context(|| format!("invalid path {}", path.display()))?;

    let (subdomain, zone_name) = args
        .domain
        .split_once('.')
        .with_context(|| format!("could not split domain {}", args.domain))?;

    let token = std::env::var("CLOUDFLARE_API_TOKEN")
        .with_context(|| "CLOUDFLARE_API_TOKEN not set")?;

    info!("Zone identifier: {zone_name}");
    info!("Subdomain: {subdomain}");

    debug!("Creating {}", path_parent.display());
    std::fs::create_dir_all(path_parent)
        .expect("Failed to create parent directories");

    let last_ip: Option<Ipv4Addr> = if path.exists() {
        fs::read_to_string(&path)
            .with_context(|| format!("could not read {}", path.display()))
            .ok()
            .and_then(|contents| contents.parse().ok())
    } else {
        None
    };

    info!(
        "Last IP address: {} ({})",
        last_ip
            .map(|value| value.to_string())
            .unwrap_or_else(|| "unknown".to_string()),
        path.display()
    );

    let ipv4 = public_ip::addr_v4()
        .await
        .with_context(|| "could not get public IPv4 address")?;

    info!("Public IPv4 address: {ipv4}");

    if last_ip != Some(ipv4) {
        let client = Client::new(
            Credentials::UserAuthToken { token },
            HttpApiClientConfig::default(),
            Environment::Production,
        )?;

        let zones = client
            .request_handle(&zone::ListZones {
                params: zone::ListZonesParams {
                    name: Some(zone_name.to_string()),
                    ..Default::default()
                },
            })
            .await?
            .result;

        if zones.len() != 1 {
            bail!("found more than one zone with name {}", zone_name);
        }

        let zone = zones
            .first()
            .with_context(|| format!("could not find zone {zone_name}"))?;

        debug!("Found zone: {zone:#?}");

        client
            .request_handle(&dns::CreateDnsRecord {
                zone_identifier: &zone.id,
                params: dns::CreateDnsRecordParams {
                    name: subdomain,
                    content: dns::DnsContent::A { content: ipv4 },
                    priority: None,
                    proxied: None,
                    ttl: None,
                },
            })
            .await
            .context("could not update DNS record")?;

        info!("Updated DNS record");

        fs::write(&path, ipv4.to_string())
            .with_context(|| format!("could not write {}", path.display()))?;

        info!("Saved last IP address to {}", path.display());
    }

    Ok(())
}
