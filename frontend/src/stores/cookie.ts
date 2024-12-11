// ตั้งค่า cookie
const setCookie = (name: string, value: string, days: number, hours: number = 0): void => {
    const d = new Date();
    if (hours > 0) {
        d.setTime(d.getTime() + (hours * 60 * 60 * 1000)); // กำหนดกี่ชั่วโมงหมดอายุ
    } else {
        d.setTime(d.getTime() + (days * 24 * 60 * 60 * 1000)); // กำหนดวันหมดอายุ
    }

    const expires = "expires=" + d.toUTCString();
    document.cookie = `${name}=${value}; ${expires}; path=/`;
};

// อ่านค่า cookie
const getCookie = (name: string): string | null => {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        const c = ca[i].trim();
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
};

const deleteCookie = (name: string): void => {
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
};

export { setCookie, getCookie, deleteCookie };