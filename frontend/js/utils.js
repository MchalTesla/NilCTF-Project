$(document).ready(function() {
    // 设置全局 AJAX 请求头
    $.ajaxSetup({
        beforeSend: function(xhr, settings) {
            const token = localStorage.getItem('token'); // 从 localStorage 获取 token
            
            if (token) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + token); // 添加 Authorization 头
            }
        }
    });
});