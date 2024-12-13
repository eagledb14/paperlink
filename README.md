# paperlink

Paperlink is a reporting and collaboration tool designed for Cyber Security Red Teams. It enables teams to document their findings during penetration tests, associate them with various assets in a client's network, and collaborate on report creation.

## Installation

1. Clone the repository from Git.
2. Ensure you have both make and npm installed on your system.

## Usage

### Running the Server

1. Navigate to the repository directory.
2. Execute the command:

```make```

This will start the server in the current directory on port 8080

### Building the Application

1. Navigate to the repository directory.
2. Run the following command:

```make build```

This will make a zip file called paperlink.zip containing all files for the application.

3. On the server you will run this on, unzip and run

```npm install```

4. Then run 

```./paperlink```

To start the application


