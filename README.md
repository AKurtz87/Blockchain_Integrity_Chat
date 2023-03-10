# Blockchain_Chat 🔗💬🔗

⚠️⚠️⚠️ **THIS REPOSITORY IN UNDER REFITTING!** ⚠️⚠️⚠️

## Description:
In the context of this chat application, integrity refers to the accuracy and consistency of the data being transmitted between clients. Data integrity is a fundamental aspect of information security, and it ensures that data is not modified, corrupted, or deleted during transmission or storage. In the case of this application, blockchain hashing is used to ensure that messages sent between clients are secure and tamper-proof. Each message is hashed using the SHA256 algorithm, and the resulting hash is added to the blockchain for verification purposes. This ensures that messages cannot be modified or deleted without detection, and it provides an additional layer of security to the application. In general, data integrity is crucial in ensuring the reliability and trustworthiness of data, and it is essential in applications where data accuracy and consistency are critical.
This Git repository contains the source code for a simple chat application built in Go that incorporates blockchain hashing for message integrity. The application allows users to send and receive messages in real-time using websockets, with each message being hashed and added to the blockchain for verification purposes.

## Features:

1. Real-time chat functionality using websockets
2. Blockchain hashing for secure message transfer
3. Simple and intuitive user interface
4. Easy to deploy and configure

## Installation and Usage:

To use this chat application, you will need to have Go installed on your system. Once you have Go installed, you can follow these steps to get the application up and running:

1. Clone the repository to your local machine using Git.
2. Navigate to the cloned repository using the command line.
3. Run the following command to install the required packages:
> go get
4. Start the chat server using the following command:
> go run main.go
5. Open your web browser and navigate to http://localhost:8080.
6. Enter your username and start chatting!

## Code Structure:

**The code is divided into four main files:**

- main.go - This file contains the main function that starts the web server and initializes the Hub.
- hub.go - This file defines the Hub struct that manages the list of connected clients and broadcasts messages to them.
- client.go - This file defines the Client struct that represents a single connected client.

## Technologies Used:

- 🔢 Go programming language
- 🦍 Gorilla WebSocket library for websockets
- 🗄️ MySQL database for storing chat history
- 🔒 SHA256 hash function for message hashing

## Contributing:

Contributions to this repository are welcome. If you notice any issues or have any suggestions for improvements, please feel free to submit a pull request or open an issue.
