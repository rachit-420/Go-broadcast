const joinPanel = document.getElementById("join-panel");
const chatPanel = document.getElementById("chat-panel");
const nameForm = document.getElementById("name-form");
const nameInput = document.getElementById("name-input");
const connectButton = document.getElementById("connect-button");
const messageForm = document.getElementById("message-form");
const messageInput = document.getElementById("message-input");
const sendButton = document.getElementById("send-button");
const leaveButton = document.getElementById("leave-button");
const messages = document.getElementById("messages");
const connectionStatus = document.getElementById("connection-status");
const youLabel = document.getElementById("you-label");

let socket;
let username = "";
let leaving = false;

nameForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const nextName = nameInput.value.trim();
    if (!nextName || socket) {
        nameInput.focus();
        return;
    }

    username = nextName;
    leaving = false;
    let connected = false;
    messages.replaceChildren();
    setStatus("Connecting…");
    connectButton.disabled = true;

    const url = new URL("ws://localhost:8080/ws");
    url.searchParams.set("username", username);
    socket = new WebSocket(url);

    socket.onopen = () => {
        connected = true;
        joinPanel.hidden = true;
        chatPanel.hidden = false;
        youLabel.textContent = username;
        messageInput.disabled = false;
        sendButton.disabled = false;
        messageInput.focus();
    };

    socket.onmessage = (event) => {
        let data;
        try {
            data = JSON.parse(event.data);
        } catch {
            appendEvent("note", "Received an unreadable message.");
            return;
        }

        if (!isChatMessage(data)) {
            appendEvent("note", "Received a message the room didn’t understand.");
            return;
        }

        renderMessage(data);
    };

    socket.onerror = () => {
        setStatus("Couldn’t connect. That name might already be taken.", true);
    };

    socket.onclose = () => {
        socket = null;
        connectButton.disabled = false;
        messageInput.disabled = true;
        sendButton.disabled = true;
        chatPanel.hidden = true;
        joinPanel.hidden = false;
        if (leaving) {
            setStatus("You left. Enter a name to go back in.");
        } else if (!connected) {
            setStatus("Couldn’t connect. That name might already be taken.", true);
        } else {
            setStatus("Disconnected. Enter a name to reconnect.", true);
        }
        leaving = false;
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

leaveButton.addEventListener("click", () => {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        return;
    }
    leaving = true;
    socket.close();
});

function isChatMessage(data) {
    return Boolean(
        data &&
        typeof data === "object" &&
        typeof data.type === "string" &&
        typeof data.userName === "string" &&
        typeof data.message === "string" &&
        typeof data.timestamp === "string"
    );
}

function renderMessage(data) {
    const type = data.type.toLowerCase();
    if (type === "joined") {
        appendEvent("joined", `${data.userName} joined`);
        return;
    }
    if (type === "left") {
        appendEvent("left", `${data.userName} left`);
        return;
    }

    appendBubble(data.userName, data.message, data.userName === username, data.timestamp);
}

function appendEvent(kind, text) {
    const item = document.createElement("li");
    item.className = `event event--${kind}`;
    item.textContent = text;
    messages.appendChild(item);
    scrollToLatest();
}

function appendBubble(name, text, mine, timestamp) {
    const item = document.createElement("li");
    item.className = mine ? "bubble-row bubble-row--mine" : "bubble-row";

    if (!mine) {
        const sender = document.createElement("span");
        sender.className = "bubble-name";
        sender.textContent = name;
        item.appendChild(sender);
    }

    const bubble = document.createElement("p");
    bubble.className = "bubble";
    bubble.textContent = text;

    const sentAt = new Date(timestamp);
    const time = document.createElement("time");
    time.className = "bubble-time";
    time.dateTime = timestamp;
    time.textContent = Number.isNaN(sentAt.getTime())
        ? ""
        : sentAt.toLocaleTimeString([], {
            hour: "numeric",
            minute: "2-digit",
        });

    const body = document.createElement("div");
    body.className = "bubble-body";
    body.append(bubble, time);
    item.append(body);
    messages.appendChild(item);
    scrollToLatest();
}

function scrollToLatest() {
    messages.scrollTop = messages.scrollHeight;
}

function setStatus(text, isError = false) {
    connectionStatus.textContent = text;
    connectionStatus.classList.toggle("status--error", isError);
}
