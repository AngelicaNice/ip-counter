import random
import os

TARGET_SIZE = 100 * 1024 * 1024 * 1024 

output_file = "ip_addresses"

def generate_random_ip():
    return f"{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}"

def generate_ips_to_file(file_name, target_size):
    with open(file_name, 'w') as file:
        current_size = 0
        while current_size < target_size:
            ip_address = generate_random_ip() + "\n"
            # Записываем IP в файл
            file.write(ip_address)
            current_size += len(ip_address)
            if current_size % (100 * 1024 * 1024) == 0:
                print(f"Current file size: {current_size / (1024 * 1024):.2f} MB")

    print(f"File '{file_name}' has reached the target size of {target_size / (1024 * 1024 * 1024):.2f} GB.")

generate_ips_to_file(output_file, TARGET_SIZE)