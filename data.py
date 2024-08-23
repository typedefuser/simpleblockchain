import json
import os
from faker import Faker
import random


fake = Faker()

def create_transaction():

    sender = fake.name()
    receiver = fake.name()

  
    amount = round(random.uniform(1, 1000), 2)

    transaction = {
        "Sender": sender,
        "Receiver": receiver,
        "Amount": amount
    }

    return transaction

def create_block(filename="block.json"):
    block=[]
    for _ in range(10):
        transaction = create_transaction()
        block.append(transaction)


    with open(filename, "w") as file:
        json.dump(block, file, indent=4)

if __name__ == "__main__":
    create_block()
