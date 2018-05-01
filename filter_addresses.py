import sys

if __name__ == '__main__':
    seen_addresses = {}

    # strip out duplicates
    with open(sys.argv[1]) as f:
        for line in f:
            line = line.strip()
            if line == '':
                break
            if not seen_addresses.get(line):
                seen_addresses[line] = True
                print(line)
