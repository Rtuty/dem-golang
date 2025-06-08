// Утилиты для удаления материалов и продуктов

function deleteMaterial(id) {
    if (confirm('Вы уверены, что хотите удалить этот материал? Это действие нельзя отменить.')) {
        fetch(`/api/v1/materials/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            // Проверяем, является ли ответ JSON
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return response.json();
            } else {
                throw new Error('Сервер вернул не JSON ответ');
            }
        })
        .then(data => {
            if (data && data.success) {
                alert('Материал успешно удален');
                location.reload();
            } else {
                alert('Ошибка: ' + (data && data.error ? data.error : 'Неизвестная ошибка'));
            }
        })
        .catch(error => {
            console.error('Ошибка удаления материала:', error);
            alert('Ошибка удаления: ' + error.message);
        });
    }
}

function deleteProduct(id) {
    if (confirm('Вы уверены, что хотите удалить эту продукцию?')) {
        fetch(`/api/v1/products/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            // Проверяем, является ли ответ JSON
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return response.json();
            } else {
                throw new Error('Сервер вернул не JSON ответ');
            }
        })
        .then(data => {
            if (data && data.success) {
                alert('Продукция успешно удалена');
                location.reload();
            } else {
                alert('Ошибка: ' + (data && data.error ? data.error : 'Неизвестная ошибка'));
            }
        })
        .catch(error => {
            console.error('Ошибка удаления продукции:', error);
            alert('Ошибка удаления: ' + error.message);
        });
    }
} 