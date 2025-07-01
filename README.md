# 🛠️ AWS S3 Tutorial

A simple CLI to **upload**, **download**, and **list files** (e.g., images) from **Amazon S3**, written in [Go](https://golang.org) using the AWS SDK v2.

---

![Go](https://img.shields.io/badge/Go-1.20%2B-blue.svg) ![AWS SDK v2](https://img.shields.io/badge/AWS%20SDK-v2-blue) ![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)

---

## 🚀 Features

-   ✅ Upload local files to S3
-   📥 Download files from S3
-   📃 List objects in a bucket
-   🔒 Uses AWS credentials via environment or `~/.aws/credentials`
-   🧱 Modular structure (`cmd/` + `storage/`)

---

## 📦 Installation

### Clone the repo

```bash
git clone https://github.com/yourusername/s3cli.git
cd s3cli
```

### Run without installing

```bash
go run . -cmd=list -bucket=my-bucket
```

### Or build the binary

```bash
go build -o s3cli .
./s3cli -cmd=upload -bucket=my-bucket -file=cat.jpg -key=images/cat.jpg
```

## 🧪 Usage

### Upload a file to S3

```bash
s3cli -cmd=upload -bucket=my-bucket -file=./certificado.png -key=certificado.png
```

### List objects in a bucket

```bash
s3cli -cmd=list -bucket=my-bucket
```

### Download a file from S3

```bash
s3cli -cmd=download -bucket=my-bucket -key=certificado.png -file=certificado-downloaded.png
```

## 🔐 AWS Credentials

### Ensure your credentials are set using one of the supported methods:

### Option 1: AWS CLI config

```bash
aws configure
```

### ption 2: Environment variables

```bash
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
export AWS_REGION=us-east-1
```
