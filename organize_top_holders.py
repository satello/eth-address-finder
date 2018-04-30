import sys
import bisect

NUMBER_OF_HOLDERS = 1000000

class Hodler:
    def __init__(self, address, balance):
        self.address = address
        self.balance = balance

    def __lt__(self, other):
        return self.balance < other.balance

    def __gt__(self, other):
        return self.balance > other.balance

    def __eq__(self, other):
        return self.balance == other.balance

    def as_list(self):
        return [self.address, self.balance]

if __name__ == '__main__':
    file_name = sys.argv[1]

    top_addresses = []

    # strip out duplicates
    with open(file_name) as f:
        for line in f:
            line = line.strip()

            if line == '':
                break

            address, balance = line.split(',')
            balance = int(balance)
            if balance > 0 and (
                len(top_addresses) < NUMBER_OF_HOLDERS or balance > top_addresses[0].balance
            ):
                if len(top_addresses) >= NUMBER_OF_HOLDERS:
                    del top_addresses[0] # remove first item in list
                hodler = Hodler(address, balance) # create new hodler
                bisect.insort(top_addresses, hodler) # insert hodler

    for holder in reversed(top_addresses):
        print("%s,%d" % (holder.address, holder.balance))
