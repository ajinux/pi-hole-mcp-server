# 🛡️ Pi-hole MCP Server

> **Monitor your home network and investigate suspicious domains with AI-powered OSINT**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)

## 🎯 Why This Project?

Built to monitor home network activity and investigate DNS queries using AI. Combines Pi-hole's DNS-level visibility with domain intelligence (DNS records, WHOIS lookups) to quickly analyze network behavior, investigate suspicious activity for security purposes, or simply understand what your devices are connecting to—all through natural language AI queries.

## ✨ Features

Pi-hole integration (network stats, client monitoring, domain analytics) + Enhanced domain intelligence (DNS records, WHOIS, OSINT prompts) + AI-powered analysis via MCP

## 🚀 Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/ajinux/pi-hole-mcp-server/main/install.sh | sh
```

<details>
<summary>Build from Source</summary>

### 1. Clone the Repository

```bash
git clone https://github.com/ajinux/pi-hole-mcp-server.git
cd pi-hole-mcp-server
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build the Server

```bash
make build
```

</details>


## ▶️ Run
```
 $ pihole-mcp
2025/12/14 17:40:39 starting pihole-mcp version=v0.1.0-dirty commit=961ceed built=2025-12-14T11:54:53Z
2025/12/14 17:40:40 obtained a new pi-hole session token successfully
2025/12/14 17:40:40 Starting Pi-hole MCP server on http://localhost:8081
2025/12/14 17:40:40 Connected to Pi-hole at: http://192.168.0.111/api
```

## ⚙️ Configuration

Create a `.env` file or set it env or pass it as arg (check pihole-mcp --help):

```env
PIHOLE_URL=http://192.168.1.100:83/api
PIHOLE_PASSWORD=your_pihole_api_password
PORT=8081 #mcp server port
```

You can use pi-phone admin dashboard to get new password
![pi-hole password](/assets/pi_hole_password.png)

**Note:** This server uses Streamable HTTP transport. Configure your MCP client accordingly (see MCP documentation for client-specific setup)

## 🔧 Available Tools

### 1. `get_top_active_clients`
Get most active devices by DNS query volume. Returns IP, name, query count, MAC address/vendor.

### 2. `get_top_domains_for_client`
Analyze DNS queries from a specific IP. Parameters: `client_ip` (required), `hours` (default: 24), `count` (default: 10).

### 3. `get_top_domains`
Get top queried domains (allowed + blocked) across all devices.

### 4. `get_domain_dns_records`
Get DNS records (A, AAAA, NS, MX, TXT) for any domain. Auto-extracts TLD from subdomains.

### 5. `get_domain_whois`
WHOIS lookup for domain registration info (registrar, dates, owner details). Auto-extracts TLD from subdomains.

## 💬 Available Prompts

### `domain-osint`
Comprehensive OSINT analysis combining DNS records, WHOIS data, hosting infrastructure, email security configs (SPF/DKIM/DMARC), and domain history to identify security concerns.

## 📚 Example Queries

- "Show me the top 10 most active devices"
- "What domains is device with ip 192.168.1.50 querying in the last 6 hours?"
- "Investigate suspicious-site.com using OSINT"
- "Get WHOIS information for example.com"

## 🙏 Acknowledgments

Thanks to [Pi-hole developers](https://pi-hole.net/) for the network-wide ad blocking solution, [likexian](https://github.com/likexian/whois) for the WHOIS package

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

[Issues](https://github.com/ajinux/pi-hole-mcp-server/issues) • [Pull Requests](https://github.com/ajinux/pi-hole-mcp-server/pulls)

---

**Made with ❤️ for home network security and privacy**
