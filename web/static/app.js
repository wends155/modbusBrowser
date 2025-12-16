const dataTableBody = document.getElementById('data-table-body');
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
        // Clear previous data
        dataTableBody.innerHTML = '';

        // Animate the table on update
        const dataTable = document.getElementById('data-table');
        dataTable.classList.add('updated');
        setTimeout(() => {
            dataTable.classList.remove('updated');
        }, 250);


        if (message.content.startsWith("Error:")) {
            const row = document.createElement('tr');
            const cell = document.createElement('td');
            cell.colSpan = 2;
            cell.textContent = message.content;
            row.appendChild(cell);
            dataTableBody.appendChild(row);
        } else {
            const pairs = message.content.split(', ');
            pairs.forEach(pair => {
                const [address, value] = pair.split(':');
                const row = document.createElement('tr');
                const addressCell = document.createElement('td');
                const valueCell = document.createElement('td');
                addressCell.textContent = address;
                valueCell.textContent = value;
                row.appendChild(addressCell);
                row.appendChild(valueCell);
                dataTableBody.appendChild(row);
            });
        }
        timestampDiv.textContent = `Last updated: ${new Date(message.timestamp).toLocaleString()}`;
    }
};

ws.onclose = () => {
    console.log('WebSocket connection closed');
    const row = document.createElement('tr');
    const cell = document.createElement('td');
    cell.colSpan = 2;
    cell.textContent = 'Connection closed. Please refresh the page.';
    row.appendChild(cell);
    dataTableBody.innerHTML = '';
    dataTableBody.appendChild(row);
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    const row = document.createElement('tr');
    const cell = document.createElement('td');
    cell.colSpan = 2;
    cell.textContent = 'An error occurred. Please check the console.';
    row.appendChild(cell);
    dataTableBody.innerHTML = '';
    dataTableBody.appendChild(row);
};