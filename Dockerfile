FROM golang:latest

# Create a directory for the app
RUN mkdir /app
 #directory to the app directory
COPY . /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build -o server . 

# Run the server executable
CMD [ "/app/server" ]



