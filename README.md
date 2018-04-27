# Ethereum Address Finder

Made with Go compiler `v1.10.1`

Simplistic script to rip through the blockchain and fetch all addresses that are a party in a transaction in a range of blocks. Requests are parallelized to speed things up. Unfortunately there is a hardcoded file descriptor limit in Geth nodes at 2048 so I limited it to sending 500 requests at once. (Can probably find a number to get closer to that limit, I kept it lower so that network activity wouldn't make the script fail).

Runs at a little over 30 blocks per second on machine with 16GB RAM and 6 CPUs.

### Usage

- You must have an Ethereum node running locally on port `8545`. You can try connecting to `infura` or some other hosted node but performance will likely suffer and they probably won't be happy with you spamming requests at them.

```
go get github.com/satello/eth-address-finder
go install
./eth-address-finder -start=<start_block> -end=<end_block> > output_file.txt
```

At the moment it just prints all addresses to `STDOUT` after it finishes so I just pipe the results into a file.
