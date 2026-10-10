const nameForm = document.getElementById("name-form");
const nameInput = document.getElementById("name-input");
const messageForm = document.getElementById("message-form");
const messageInput = document.getElementById("message-input");
const sendButton = document.getElementById("send-button");
const messages = document.getElementById("messages");
const connectionStatus = document.getElementById("connection-status");
let socket;

nameForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const username = nameInput.value.trim();
    if (!username) {
        nameInput.focus();
        return;
    }

    const url = new URL("ws://localhost:8080/ws");
    url.searchParams.set("username", username);
    socket = new WebSocket(url);
    nameForm.hidden = true;
    connectionStatus.textContent = "Connecting…";

    socket.onopen = () => {
        connectionStatus.textContent = `Connected as ${username}`;
        messageForm.hidden = false;
        messageInput.disabled = false;
        sendButton.disabled = false;
        messageInput.focus();
    };

    socket.onmessage = (event) => {
        let data;
        try {
            data = JSON.parse(event.data);
        } catch {
            connectionStatus.textContent = "Received an invalid message from the server.";
            return;
        }

        if (
            !data ||
            typeof data !== "object" ||
            typeof data.type !== "string" ||
            typeof data.userName !== "string" ||
            typeof data.message !== "string"
        ) {
            connectionStatus.textContent = "Received a message with an invalid format.";
            return;
        }

        const item = document.createElement("li");
        const sender = document.createElement("span");
        sender.textContent = data.userName;

        const type = data.type.toLowerCase();
        if (type === "joined") {
            sender.classList.add("message-username--joined");
        } else if (type === "leave") {
            sender.classList.add("message-username--leave");
        }

        item.append(sender, document.createTextNode(`: ${data.message}`));
        messages.appendChild(item);
        messages.scrollTop = messages.scrollHeight;
    };

    socket.onerror = () => {
        connectionStatus.textContent = "Connection error";
    };

    socket.onclose = () => {
        connectionStatus.textContent = "Disconnected. Enter your name to reconnect.";
        messageForm.hidden = true;
        messageInput.disabled = true;
        sendButton.disabled = true;
        nameForm.hidden = false;
        nameInput.focus();
    };
});

messageForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const message = messageInput.value.trim();
    if (!message || !socket || socket.readyState !== WebSocket.OPEN) {
        return;
    }

    socket.send(message);
    messageInput.value = "";
    messageInput.focus();
});
