

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

export function GetTimestamp(dateString: string): string {
    const d = new Date(dateString);
    return `${GetDayFullDate(dateString)} @ ${d.toLocaleTimeString().replace(/(.*)\D\d+/, '$1')}`;
}

export function GetAddressLink(a: Common.Address) {
    if (a && a.street) {
        return 'https://maps.google.com/?q=' + encodeURIComponent(a.apt + ' ' + a.street + ', ' + a.city + ', ' + a.state + ' ' + a.zip);
    }
    return '';
}

export function GetAddress(a: Common.Address) {
    if (a && a.street) {
        let apt = '';
        if (a.apt !== undefined && a.apt !== '') {
            apt = '#' + a.apt + ' ';
        }
        return apt + a.street + ', ' + a.city;
    }
    return '';
}
