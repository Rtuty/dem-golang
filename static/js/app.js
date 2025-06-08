// Основной JavaScript файл для приложения
document.addEventListener('DOMContentLoaded', function() {
    console.log('Приложение "Наш декор" загружено');
    
    // Подсветка активной навигации
    highlightActiveNav();
    
    // Инициализация интерактивных элементов
    initInteractiveElements();
});

// Подсветка активной ссылки в навигации
function highlightActiveNav() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-link');
    
    navLinks.forEach(link => {
        const href = link.getAttribute('href');
        if (currentPath === href || (href !== '/' && currentPath.startsWith(href))) {
            link.style.backgroundColor = 'rgba(255,255,255,0.3)';
        }
    });
}

// Инициализация интерактивных элементов
function initInteractiveElements() {
    // Форматирование чисел в таблицах
    formatNumbers();
    
    // Подтверждение удаления
    initDeleteConfirmations();
    
    // Валидация форм
    initFormValidation();
}

// Форматирование чисел с разделителями тысяч
function formatNumbers() {
    const priceElements = document.querySelectorAll('.product-price, .product-calculated-price, .material-cost, .material-total');
    
    priceElements.forEach(element => {
        const text = element.textContent.trim();
        const number = parseFloat(text);
        
        if (!isNaN(number)) {
            element.textContent = number.toLocaleString('ru-RU', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            });
        }
    });
}

// Инициализация подтверждений удаления
function initDeleteConfirmations() {
    const deleteButtons = document.querySelectorAll('[onclick*="deleteProduct"]');
    
    deleteButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const productId = this.getAttribute('onclick').match(/\d+/)[0];
            confirmDelete(productId);
        });
    });
}

// Подтверждение удаления продукции
function confirmDelete(productId) {
    if (confirm('Вы уверены, что хотите удалить эту продукцию? Это действие нельзя отменить.')) {
        deleteProduct(productId);
    }
}

// Удаление продукции через API
function deleteProduct(id) {
    showLoading(true);
    
    fetch(`/api/v1/products/${id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        return response.json();
    })
    .then(data => {
        showLoading(false);
        
        if (data.success || data.message) {
            showNotification('Продукция успешно удалена', 'success');
            // Удаляем строку из таблицы
            const row = document.querySelector(`tr[data-id="${id}"]`);
            if (row) {
                row.remove();
            }
            
            // Если больше нет продукции, показываем пустое состояние
            const tbody = document.querySelector('.products-table tbody');
            if (tbody && tbody.children.length === 0) {
                location.reload();
            }
        } else {
            showNotification('Ошибка: ' + (data.error || 'Неизвестная ошибка'), 'error');
        }
    })
    .catch(error => {
        showLoading(false);
        showNotification('Ошибка удаления: ' + error.message, 'error');
    });
}

// Удаление материала через API
function deleteMaterial(id) {
    if (!confirm('Вы уверены, что хотите удалить этот материал?')) {
        return;
    }
    
    showLoading(true);
    
    fetch(`/api/v1/materials/${id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        return response.json();
    })
    .then(data => {
        showLoading(false);
        
        if (data.success || data.message) {
            showNotification('Материал успешно удален', 'success');
            // Удаляем строку из таблицы
            const row = document.querySelector(`tr[data-id="${id}"]`);
            if (row) {
                row.remove();
            }
            
            // Если больше нет материалов, показываем пустое состояние
            const tbody = document.querySelector('.materials-table tbody');
            if (tbody && tbody.children.length === 0) {
                location.reload();
            }
        } else {
            showNotification('Ошибка: ' + (data.error || 'Неизвестная ошибка'), 'error');
        }
    })
    .catch(error => {
        showLoading(false);
        showNotification('Ошибка удаления: ' + error.message, 'error');
    });
}

// Показать/скрыть индикатор загрузки
function showLoading(show) {
    let loader = document.getElementById('loading-indicator');
    
    if (show && !loader) {
        loader = document.createElement('div');
        loader.id = 'loading-indicator';
        loader.innerHTML = `
            <div style="
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: rgba(0,0,0,0.5);
                display: flex;
                justify-content: center;
                align-items: center;
                z-index: 9999;
            ">
                <div style="
                    background: white;
                    padding: 2rem;
                    border-radius: 8px;
                    text-align: center;
                ">
                    <div style="
                        border: 4px solid #f3f3f3;
                        border-top: 4px solid #667eea;
                        border-radius: 50%;
                        width: 50px;
                        height: 50px;
                        animation: spin 1s linear infinite;
                        margin: 0 auto 1rem;
                    "></div>
                    <p>Загрузка...</p>
                </div>
            </div>
        `;
        
        // Добавляем CSS анимацию
        if (!document.getElementById('loading-styles')) {
            const style = document.createElement('style');
            style.id = 'loading-styles';
            style.textContent = `
                @keyframes spin {
                    0% { transform: rotate(0deg); }
                    100% { transform: rotate(360deg); }
                }
            `;
            document.head.appendChild(style);
        }
        
        document.body.appendChild(loader);
    } else if (!show && loader) {
        loader.remove();
    }
}

// Показать уведомление
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div style="
            position: fixed;
            top: 20px;
            right: 20px;
            background: ${type === 'success' ? '#d4edda' : '#f8d7da'};
            color: ${type === 'success' ? '#155724' : '#721c24'};
            padding: 1rem 1.5rem;
            border-radius: 8px;
            border-left: 4px solid ${type === 'success' ? '#28a745' : '#dc3545'};
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            z-index: 10000;
            max-width: 400px;
            word-wrap: break-word;
        ">
            ${message}
            <button onclick="this.parentElement.remove()" style="
                float: right;
                background: none;
                border: none;
                font-size: 18px;
                cursor: pointer;
                margin-left: 1rem;
            ">&times;</button>
        </div>
    `;
    
    document.body.appendChild(notification);
    
    // Автоматически скрыть через 5 секунд
    setTimeout(() => {
        if (notification.parentElement) {
            notification.remove();
        }
    }, 5000);
}

// Инициализация валидации форм
function initFormValidation() {
    const forms = document.querySelectorAll('form');
    
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            if (!validateForm(this)) {
                e.preventDefault();
            }
        });
        
        // Валидация в реальном времени
        const inputs = form.querySelectorAll('input[type="number"]');
        inputs.forEach(input => {
            input.addEventListener('input', function() {
                validateNumberInput(this);
            });
        });
    });
}

// Валидация формы
function validateForm(form) {
    let isValid = true;
    const errors = [];
    
    // Проверка обязательных полей
    const requiredFields = form.querySelectorAll('[required]');
    requiredFields.forEach(field => {
        if (!field.value.trim()) {
            isValid = false;
            errors.push(`Поле "${getFieldLabel(field)}" обязательно для заполнения`);
            highlightError(field);
        } else {
            clearError(field);
        }
    });
    
    // Проверка числовых полей
    const numberFields = form.querySelectorAll('input[type="number"]');
    numberFields.forEach(field => {
        if (field.value && !validateNumberInput(field)) {
            isValid = false;
        }
    });
    
    if (!isValid && errors.length > 0) {
        showNotification(errors.join('\n'), 'error');
    }
    
    return isValid;
}

// Валидация числового поля
function validateNumberInput(input) {
    const value = parseFloat(input.value);
    const min = parseFloat(input.getAttribute('min'));
    const max = parseFloat(input.getAttribute('max'));
    
    let isValid = true;
    let errorMessage = '';
    
    if (input.value && isNaN(value)) {
        isValid = false;
        errorMessage = 'Введите корректное число';
    } else if (!isNaN(min) && value < min) {
        isValid = false;
        errorMessage = `Значение должно быть не менее ${min}`;
    } else if (!isNaN(max) && value > max) {
        isValid = false;
        errorMessage = `Значение должно быть не более ${max}`;
    }
    
    if (isValid) {
        clearError(input);
    } else {
        highlightError(input, errorMessage);
    }
    
    return isValid;
}

// Подсветка ошибки поля
function highlightError(field, message = '') {
    field.style.borderColor = '#dc3545';
    field.style.backgroundColor = '#fff5f5';
    
    // Удаляем предыдущее сообщение об ошибке
    const existingError = field.parentElement.querySelector('.error-message');
    if (existingError) {
        existingError.remove();
    }
    
    // Добавляем новое сообщение об ошибке
    if (message) {
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.style.color = '#dc3545';
        errorDiv.style.fontSize = '12px';
        errorDiv.style.marginTop = '0.25rem';
        errorDiv.textContent = message;
        field.parentElement.appendChild(errorDiv);
    }
}

// Очистка ошибки поля
function clearError(field) {
    field.style.borderColor = '';
    field.style.backgroundColor = '';
    
    const errorMessage = field.parentElement.querySelector('.error-message');
    if (errorMessage) {
        errorMessage.remove();
    }
}

// Получение подписи поля
function getFieldLabel(field) {
    const label = field.parentElement.querySelector('label');
    return label ? label.textContent.replace('*', '').trim() : field.name;
}

// Утилита для копирования в буфер обмена
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showNotification('Скопировано в буфер обмена', 'success');
    }).catch(err => {
        console.error('Ошибка копирования: ', err);
        showNotification('Ошибка копирования', 'error');
    });
}

// Экспорт функций для глобального использования
window.deleteProduct = deleteProduct;
window.showNotification = showNotification;
window.copyToClipboard = copyToClipboard; 