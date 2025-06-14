<p align="center">
    <img src="https://github.com/PumpkinSeed/letherscan/blob/cb354c7b4909013695e35f61f2c535794b257a34/assets/logo.png" width="250" alt="logo">
</p>

# letherscan - Lightweight Etherscan

## Motivation

I tried to find an Etherscan alternative what I can run and quickly check my hardhat node's transactions. Since I couldn't find one, I decided to build it myself. The goal is to have a lightweight Etherscan alternative that can be run locally and is easy to use.

## Disclaimer

I'm a backend engineer, so the content of the frontend folder is 100% Cursor generated code. Even the logo is generated by ChatGPT.

Also, primarily I use this for my own purposes, so the code is not production ready. If you want to use it, feel free to do so, but be aware that it might not work as expected. I'm happy to accept issues and PRs, but I can't guarantee that I will fix them in a timely manner.

I built it in a night, so it's not a polished product. :D

## Build & Run - Binary

### Build the frontend

```bash
npm --prefix ./frontend run build
```

### Move the build to the backend

```bash
rm -rf ./bin/build
mv ./frontend/build ./bin
```

### Build the backend

```bash
go build -o ./bin/letherscan ./bin/main.go

./bin/letherscan
```

## Run - Docker

```bash
docker run --network=host pumpkinseed/letherscan:v0.0.5
```

## Build & Run - Docker

### Build the Docker image

```bash
docker build -t letherscan .
```

### Run the Docker container

```bash
# Host network is required for the application to access the locally running hardhat node.
docker run --network=host letherscan
```

## Screenshots

<img src="https://github.com/PumpkinSeed/letherscan/blob/cb354c7b4909013695e35f61f2c535794b257a34/assets/block_view.png" alt="block view">
<img src="https://github.com/PumpkinSeed/letherscan/blob/cb354c7b4909013695e35f61f2c535794b257a34/assets/transaction_view.png" alt="transaction view">