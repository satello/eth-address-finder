import csv

if __name__ == '__main__':
    seen_addresses = {}

    # strip out duplicates
    with open('./result.csv') as f:
        for line in f:
            line = line.strip()
            if line == '':
                break
            if not seen_addresses.get(line):
                seen_addresses[line] = True
                print(line)
