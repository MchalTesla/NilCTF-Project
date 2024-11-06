// 计算 SHA-256 哈希值的函数
async function sha256(message) {
    // 将字符串转换为 ArrayBuffer
    const msgBuffer = new TextEncoder().encode(message);
    
    // 使用 SubtleCrypto API 计算 SHA-256 哈希
    const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
    
    // 将 ArrayBuffer 转换为十六进制字符串
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    
    return hashHex;
}
