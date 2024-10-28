$(document).ready(function() {
    // 注册表单提交处理
    if ($('#register-form').length) {
        $('#register-form').submit(async function(event) {
            event.preventDefault(); // 阻止默认表单提交

            const username = $('#register-username').val();
            const password = $('#register-password').val();
            const email = $('#register-email').val(); // 新增邮箱字段

            // 散列密码
            const hashedPassword = await hashPassword(password);

            $.ajax({
                type: 'POST',
                url: '/api/register', // 你的注册API端点
                contentType: 'application/json',
                data: JSON.stringify({ username, password: hashedPassword, email }), // 包含邮箱字段
                success: function(response) {
                    alert(response.message); // 显示成功消息
                    if (response.redirect) {
                        window.location.href = response.redirect;
                    }
                    $('#register-form')[0].reset();
                },
                error: function(xhr) {
                    alert(xhr.responseJSON.message); // 显示错误消息
                }
            });
        });
    }

    // 登录表单提交处理
    if ($('#login-form').length) {
        $('#login-form').submit(async function(event) {
            event.preventDefault(); // 阻止默认表单提交

            const username = $('#login-username').val();
            const password = $('#login-password').val();

            // 散列密码
            const hashedPassword = await hashPassword(password);

            $.ajax({
                type: 'POST',
                url: '/api/login', // 你的登录API端点
                contentType: 'application/json',
                data: JSON.stringify({ username, password: hashedPassword }),
                success: function(response) {
                    // 保存 token 到 localStorage
                    if (response.token) {
                        localStorage.setItem('token', response.token);
                    }
                    // 登录成功后重定向到 /index
                    if (response.redirect) {
                        window.location.href = response.redirect;
                    }
                },
                error: function(xhr) {
                    alert(xhr.responseJSON.message); // 显示错误消息
                }
            });
        });
    }

    // 密码散列函数
    async function hashPassword(password) {
        const encoder = new TextEncoder();
        const data = encoder.encode(password);
        const hash = await crypto.subtle.digest('SHA-256', data); // 使用SHA-256算法进行散列
        const hashArray = Array.from(new Uint8Array(hash)); // 将Hash值转换为字节数组
        const hashedPassword = hashArray.map(b => ('00' + b.toString(16)).slice(-2)).join(''); // 转换为十六进制字符串
        return hashedPassword;
    }
});