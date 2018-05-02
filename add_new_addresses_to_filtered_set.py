import sys

if __name__ == '__main__':
    seen_addresses = {}
    new_addresses = []

    extra_address_file = sys.argv[1]
    # add all of our new addresses (smaller than current address set)
    with open(extra_address_file) as f:
        for line in f:
            line = line.strip()
            if line == '':
                break
            if not seen_addresses.get(line):
                try:
                    seen_addresses[line] = True
                except:
                    raise RuntimeError('%s made a memory boo boo' % line)
                new_addresses.append(line)

    filtered_file = sys.argv[2]
    with open(filtered_file) as f:
        for line in f:
            line = line.strip()
            if line == '':
                break
            if seen_addresses.get(line):
                # take out address from array
                new_addresses.remove(line)
                # mark as false so we just skip over it next time
                seen_addresses[line] = False
                print(line)

    # print out our new addresses
    for new_address in new_addresses:
        print(new_address)
