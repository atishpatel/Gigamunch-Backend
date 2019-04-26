

// Monday, January 1
export function GetDayMonthDayDate(dateString: string): string {
    const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wedensday', 'Thursday', 'Friday', 'Saturday'];
    const d = new Date(dateString);
    const day = d.getUTCDay();
    const monthDayDate = GetMonthDayDate(dateString);
    return `${dayNames[day]}, ${monthDayDate}`;
}

// January 1
export function GetMonthDayDate(dateString: string): string {
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    const d = new Date(dateString);
    const month = d.getUTCMonth();
    const date = d.getUTCDate();
    return `${monthNames[month]} ${date}`;
}


// Monday, January 1, 2018
export function GetDayFullDate(dateString: string): string {
    const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wedensday', 'Thursday', 'Friday', 'Saturday'];
    const d = new Date(dateString);
    const day = d.getUTCDay();
    const fullDate = GetFullDate(dateString);
    return `${dayNames[day]}, ${fullDate}`;
}

// January 1, 2018
export function GetFullDate(dateString: string): string {
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    const d = new Date(dateString);
    const month = d.getUTCMonth();
    const date = d.getUTCDate();
    const year = d.getUTCFullYear();
    return `${monthNames[month]} ${date}, ${year}`;
}
