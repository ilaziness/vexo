/**
 * 将任意支持的数据编码为 Base64 字符串
 * 支持: string | Uint8Array | ArrayBuffer | Buffer（在 Node.js 中）
 */
export function encodeBase64(data: string | Uint8Array | ArrayBuffer): string {
    let uint8Array: Uint8Array;

    if (typeof data === 'string') {
        // 字符串 → UTF-8 编码的 Uint8Array
        const encoder = new TextEncoder();
        uint8Array = encoder.encode(data);
    } else if (data instanceof Uint8Array) {
        uint8Array = data;
    } else if (data instanceof ArrayBuffer) {
        uint8Array = new Uint8Array(data);
    } else {
        // 如果未来传入其他类型（如 Buffer），可尝试 fallback
        // 注意：Buffer 在浏览器中不存在，此处不显式处理
        throw new TypeError('Unsupported data type for encodeBase64');
    }

    // 将 Uint8Array 转为 Latin1 字符串，供 btoa 使用
    let binary = '';
    const len = uint8Array.length;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(uint8Array[i]);
    }

    return btoa(binary);
}

/**
 * 将 Base64 字符串解码为 Uint8Array（原始二进制数据）
 * 自动忽略空白字符（如换行、空格）
 */
export function decodeBase64(base64Str: string): Uint8Array {
    // 移除所有空白字符（兼容 MIME/Base64 with newlines）
    const cleaned = base64Str.replace(/\s/g, '');

    let binary: string;
    try {
        binary = atob(cleaned);
    } catch (e) {
        throw new Error(`Invalid base64 string: ${e instanceof Error ? e.message : 'Unknown error'}`);
    }

    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
}
