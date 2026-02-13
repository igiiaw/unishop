// Cart functionality for updating quantities
function updateQuantity(itemId, change) {
    const qtyInput = document.getElementById('qty-' + itemId);
    let currentQty = parseInt(qtyInput.value);
    let newQty = currentQty + change;

    // Ensure quantity is at least 1
    if (newQty < 1) {
        newQty = 1;
    }

    // Update the input field
    qtyInput.value = newQty;

    // Send update to server
    fetch('/cart/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'item_id=' + itemId + '&quantity=' + newQty
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'ok') {
                // Reload page to update totals
                location.reload();
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to update quantity');
        });
}