const dataDiv = document.getElementById('data');
const serverInfoDiv = document.getElementById('server-info');
const timestampDiv = document.getElementById('timestamp');
const ws = new WebSocket(`ws://${window.location.host}/ws`);

ws.onopen = () => {
    console.log('WebSocket connection opened');
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'serverInfo') {
        serverInfoDiv.textContent = message.content;
    } else if (message.type === 'modbusData') {
        dataDiv.textContent = message.content;
        timestampDiv.textContent = `Last updated: ${new Date(message.timestamp).toLocaleString()}`;
    }
};

ws.onclose = () => {
    console.log('WebSocket connection closed');
    dataDiv.textContent = 'Connection closed. Please refresh the page.';
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    dataDiv.textContent = 'An error occurred. Please check the console.';
};