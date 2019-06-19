import { GetAddressLink, GetAddress, GetDayFullDate, GetFullDate, GetTimestamp } from '../ts/utils';


export function GetSubscribersExtended(s: Types.SubscriberExtended[] | Common.Subscriber[]): Types.SubscriberExtended[] {
    const subs = s as Types.SubscriberExtended[];
    for (let i = 0; i < subs.length; i++) {
        subs[i] = GetSubscriberExtended(subs[i]);
    }
    return subs;
}

export function GetSubscriberExtended(s: Types.SubscriberExtended | Common.Subscriber): Types.SubscriberExtended {
    const sub = s as Types.SubscriberExtended;
    sub.addressString = GetAddress(sub.address);
    sub.addressLink = GetAddressLink(sub.address);
    sub.emails = [];
    sub.names = [];
    sub.email_prefs.reduce((emails, emailPrefs, a, b) => {
        emails.push(emailPrefs.email);
        sub.names.push(
            `${emailPrefs.first_name} ${emailPrefs.last_name}`
        );
        return emails;
    }, sub.emails);
    sub.emailsString = sub.emails.toString();
    sub.namesString = sub.names.toString();
    sub.phonenumbers = [];
    sub.phone_prefs.reduce((numbers, phonePrefs, a, b) => {
        numbers.push(phonePrefs.number);
        return numbers;
    }, sub.phonenumbers);
    sub.phonenumbersString = sub.phonenumbers.toString();
    sub.signUpDatetimeTimestamp = GetTimestamp(sub.sign_up_datetime);
    sub.activateDatetimeTimestamp = GetTimestamp(sub.activate_datetime);
    if (sub.deactivatedDatetimeTimestamp) {
        sub.deactivatedDatetimeTimestamp = GetTimestamp(sub.deactivated_datetime);
    } else {
        sub.deactivatedDatetimeTimestamp = '-';
    }
    return sub;
}


export function GetActivitiesExtended(v: Types.ActivityExtended[] | Common.Activity[]): Types.ActivityExtended[] {
    const acts = v as Types.ActivityExtended[];
    for (let i = 0; i < acts.length; i++) {
        acts[i] = GetActivityExtended(acts[i]);
    }
    return acts;
}

export function GetActivityExtended(v: Types.ActivityExtended | Common.Activity): Types.ActivityExtended {
    const act = v as Types.ActivityExtended;
    act.dateFull = GetDayFullDate(act.date);
    act.status = getActivityStatus(act);
    act.addressString = act.address_string;
    act.discountString = getActivityDiscountString(act.discount_amount, act.discount_percent);
    act.paidDate = GetFullDate(act.paid_datetime);
    return act;
}

export function GetLogsExtended(v: Types.LogExtended[] | Common.Log[]): Types.LogExtended[] {
    const logs = v as Types.LogExtended[];
    for (let i = 0; i < logs.length; i++) {
        logs[i] = GetLogExtended(logs[i]);
    }
    return logs;
}

export function GetLogExtended(v: Types.LogExtended | Common.Log): Types.LogExtended {
    const log = v as Types.LogExtended;
    log.color = getLogColor(log);
    log.icon = getLogIcon(log);
    log.timestampString = GetTimestamp(log.timestamp);
    log.basicPayloadDescriptionHTML = log.basic_payload.description.replace(/;;;/g, '<br>');
    return log;
}

export function GetUnpaidSummariesExtended(v: Types.UnpaidSummaryExtended[] | AdminAPI.UnpaidSummary[]): Types.UnpaidSummaryExtended[] {
    const ves = v as Types.UnpaidSummaryExtended[];
    for (let i = 0; i < ves.length; i++) {
        ves[i] = GetUnpaidSummaryExtended(ves[i]);
    }
    return ves;
}

export function GetUnpaidSummaryExtended(v: Types.UnpaidSummaryExtended | AdminAPI.UnpaidSummary): Types.UnpaidSummaryExtended {
    const ve = v as Types.UnpaidSummaryExtended;
    ve.name = `${v.first_name} ${v.last_name}`;
    return ve;
}

function getLogColor(log: Types.LogExtended) {
    const type = log.type;
    const action = log.action;
    switch (action) {
        case 'skip':
            return 'orange';
        case 'unskip':
            return 'pink';
        case 'message':
            if (log.basic_payload.title.includes('from(Gigamunch)')) {
                return 'green';
            }
            return 'light-green';
        case 'rating':
            return 'cyan';
        case 'update':
        case 'servings_changed':
        case 'servings_changed_permanently':
        case 'card_updated':
            return 'amber';
        case 'refund':
            return 'red';
    }
    return 'blue-grey';
}

function getLogIcon(log: Types.LogExtended) {
    const type = log.type;
    const action = log.action;
    switch (action) {
        case 'skip':
            return 'remove_shopping_cart';
        case 'unskip':
            return 'add_shopping_cart';
        case 'message':
            if (log.basic_payload.title.includes('from(Gigamunch)')) {
                return 'reply';
            }
            return 'message';
        case 'rating':
            return 'star_rate';
        case 'update':
        case 'servings_changed':
        case 'servings_changed_permanently':
        case 'card_updated':
            return 'cloud_upload';
        case 'paid':
            return 'attach_money';
        case 'refund':
            return 'money_off';
    }
    return 'bubble_chart';
}


function getActivityStatus(act: Types.ActivityExtended) {
    if (act.refunded) {
        return 'Refunded $' + act.refunded_amount;
    } else if (act.skip) {
        return 'Skipped';
    } else if (act.first) {
        return 'First';
    } else if (act.paid) {
        return 'Paid $' + act.amount_paid;
    } else if (act.discount_percent === 100) {
        return `Discounted ${act.discount_percent}%`;
    }
    const today = new Date();
    const d = new Date(act.date);
    if (today < d) {
        return 'Pending';
    }
    return 'Owe $' + act.amount;
}

function getActivityDiscountString(discountAmount: number, discountPercent: number) {
    if (discountAmount > 0 && discountPercent > 0) {
        return `$${discountAmount} | ${discountPercent}%`;
    } else if (discountAmount > 0) {
        return `$${discountAmount}`;
    } else if (discountPercent > 0) {
        return `${discountPercent}%`;
    } else {
        return '–';
    }
}
