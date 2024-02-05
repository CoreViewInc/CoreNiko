# Use the official Ubuntu base image
FROM ubuntu:latest

# Set the working directory in the container
WORKDIR /usr/src/app

# Install Python 3
RUN apt-get update && \
    apt-get install -y python3

# Copy the current directory contents into the container at /usr/src/app
COPY . .

# Run a Python script when the container launches
CMD ["python3", "script.py"]