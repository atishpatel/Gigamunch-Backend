COOK = COOK || {};
class ServiceOld {
    constructor() {
        this.loaded = false;
        this.callQueue = [];
    }
    endpointLoaded() {
        this.loaded = true;
        this.service = gapi.client.cookservice;
        if (this.callQueue) {
            for (const fn of this.callQueue) {
                fn();
            }
            this.callQueue = [];
        }
        setTimeout(() => {
            this.refreshToken();
        }, 3000);
        ga('send', {
            hitType: 'timing',
            timingCategory: 'endpoint',
            timingVar: 'load',
            timingValue: window.performance.now(),
        });
    }
    schedulePhoneCall(phoneNumber, datetime, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.schedulePhoneCall(phoneNumber, datetime, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            datetime,
            phone_number: phoneNumber,
        };
        this
            .service
            .schedulePhoneCall(request)
            .execute((resp) => {
            this.logError('schedulePhoneCall', resp.err);
            callback(resp.err);
        });
    }
    finishOnboarding(cook, submerchant, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.finishOnboarding(cook, submerchant, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            cook,
            sub_merchant: submerchant,
        };
        this
            .service
            .finishOnboarding(request)
            .execute((resp) => {
            this.logError('finishOnboarding', resp.err);
            COOK.User.update(resp.gigatoken);
            callback(resp.err);
        });
    }
    getCook(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getCook(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getCook(request)
            .execute((resp) => {
            this.logError('getCook', resp.err);
            callback(resp.cook, resp.err);
        });
    }
    updateCook(cook, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.updateCook(cook, callback);
            });
            return;
        }
        cook.delivery_price = String(cook.delivery_price);
        const request = {
            gigatoken: this.getToken(),
            cook,
        };
        this
            .service
            .updateCook(request)
            .execute((resp) => {
            this.logError('updateCook', resp.err);
            callback(resp.cook, resp.err);
            setTimeout(() => {
                this.refreshToken();
            }, 1);
        });
    }
    getSubMerchant(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getSubMerchant(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubMerchant(request)
            .execute((resp) => {
            this.logError('getSubMerchant', resp.err);
            callback(resp.sub_merchant, resp.err);
        });
    }
    updateSubMerchant(submerchant, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.updateSubMerchant(submerchant, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            sub_merchant: submerchant,
        };
        this
            .service
            .updateSubMerchant(request)
            .execute((resp) => {
            this.logError('updateSubMerchant', resp.err);
            callback(resp.cook, resp.err);
            setTimeout(() => {
                this.refreshToken();
            }, 1);
        });
    }
    getMessageToken(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getMessageToken(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            device_id: 'browser',
        };
        this
            .service
            .getMessageToken(request)
            .execute((resp) => {
            this.logError('getMessageToken', resp.err);
            callback(resp.token, resp.err);
        });
    }
    getInquiry(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getInquiry(id, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            id,
        };
        this
            .service
            .getInquiry(request)
            .execute((resp) => {
            this.logError('getInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    }
    getInquiries(startIndex, endIndex, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getInquiries(startIndex, endIndex, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            start_index: startIndex,
            end_index: endIndex,
        };
        this
            .service
            .getInquiries(request)
            .execute((resp) => {
            this.logError('getInquiries', resp.err);
            callback(resp.inquiries, resp.err);
        });
    }
    acceptInquiry(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.acceptInquiry(id, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            id,
        };
        this
            .service
            .acceptInquiry(request)
            .execute((resp) => {
            this.logError('acceptInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    }
    declineInquiry(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.declineInquiry(id, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            id,
        };
        this
            .service
            .declineInquiry(request)
            .execute((resp) => {
            this.logError('declineInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    }
    getMenus(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getMenus(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getMenus(request)
            .execute((resp) => {
            this.logError('getMenus', resp.err);
            callback(resp.menus, resp.err);
        });
    }
    saveMenu(menu, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.saveMenu(menu, callback);
            });
            return;
        }
        delete menu.items;
        const request = {
            gigatoken: this.getToken(),
            menu,
        };
        this
            .service
            .saveMenu(request)
            .execute((resp) => {
            this.logError('saveMenu', resp.err);
            callback(resp.menu, resp.err);
        });
    }
    getItem(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getItem(id, callback);
            });
            return;
        }
        const request = {
            id,
            gigatoken: this.getToken(),
        };
        this
            .service
            .getItem(request)
            .execute((resp) => {
            this.logError('getItem', resp.err);
            callback(resp.item, resp.err);
        });
    }
    saveItem(item, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.saveItem(item, callback);
            });
            return;
        }
        item.min_servings = Number(item.min_servings);
        item.max_servings = Number(item.max_servings);
        item.cook_price_per_serving = Number(item.cook_price_per_serving);
        const request = {
            gigatoken: this.getToken(),
            item,
        };
        this
            .service
            .saveItem(request)
            .execute((resp) => {
            this.logError('saveItem', resp.err);
            callback(resp.item, resp.err);
        });
    }
    activateItem(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.activateItem(id, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            id,
        };
        this
            .service
            .activateItem(request)
            .execute((resp) => {
            this.logError('activateItem', resp.err);
            callback(resp.err);
        });
    }
    deactivateItem(id, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.deactivateItem(id, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            id,
        };
        this
            .service
            .deactivateItem(request)
            .execute((resp) => {
            this.logError('deactivateItem', resp.err);
            callback(resp.err);
        });
    }
    getSubLogs(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getSubLogs(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubLogs(request)
            .execute((resp) => {
            this.logError('getSubLogs', resp.err);
            callback(resp.sublogs, resp.err);
        });
    }
    getSubLogsForDate(date, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getSubLogsForDate(date, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
        };
        this
            .service
            .getSubLogsForDate(request)
            .execute((resp) => {
            this.logError('getSubLogsForDate', resp.err);
            callback(resp.sublogs, resp.err);
        });
    }
    getSubEmails(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getSubEmails(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubEmails(request)
            .execute((resp) => {
            this.logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.err);
        });
    }
    getSubEmailsAndSubs(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.getSubEmailsAndSubs(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubEmails(request)
            .execute((resp) => {
            this.logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.subscribers, resp.err);
        });
    }
    skipSubLog(date, subEmail, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.skipSubLog(date, subEmail, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
            sub_email: subEmail,
        };
        this
            .service
            .skipSubLog(request)
            .execute((resp) => {
            this.logError('skipSubLog', resp.err);
            callback(resp.err);
        });
    }
    CancelSub(email, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.CancelSub(email, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            email: email,
        };
        this
            .service
            .CancelSub(request)
            .execute((resp) => {
            this.logError('CancelSub', resp.err);
            callback(resp.err);
        });
    }
    discountSubLog(date, subEmail, amount, percent, overrideDiscount, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.discountSubLog(date, subEmail, amount, percent, overrideDiscount, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
            sub_email: subEmail,
            amount: amount,
            percent: percent,
            override_discount: overrideDiscount,
        };
        this
            .service
            .DiscountSubLog(request)
            .execute((resp) => {
            this.logError('DiscountSubLog', resp.err);
            callback(resp.err);
        });
    }
    ChangeServingsForDate(date, subEmail, servings, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.ChangeServingsForDate(date, subEmail, servings, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
            sub_email: subEmail,
            servings: servings,
        };
        this
            .service
            .ChangeServingsForDate(request)
            .execute((resp) => {
            this.logError('ChangeServingForDate', resp.err);
            callback(resp.err);
        });
    }
    ChangeServingsPermanently(email, servings, vegetarian, callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.ChangeServingsPermanently(email, servings, vegetarian, callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
            email: email,
            servings: servings,
            vegetarian: vegetarian,
        };
        this
            .service
            .ChangeServingsPermanently(request)
            .execute((resp) => {
            this.logError('ChangeServingsPermanently', resp.err);
            callback(resp.err);
        });
    }
    GetGeneralStats(callback) {
        if (!this.loaded) {
            this.callQueue.push(() => {
                this.GetGeneralStats(callback);
            });
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .GetGeneralStats(request)
            .execute((resp) => {
            this.logError('GetGeneralStats', resp.err);
            callback(resp);
        });
    }
    refreshToken() {
        if (!this.loaded) {
            this.callQueue.push(this.refreshToken);
            return;
        }
        const request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .refreshToken(request)
            .execute((resp) => {
            if (!this.logError('refreshToken', resp.err)) {
                COOK.User.update(resp.gigatoken);
            }
        });
    }
    getToken() {
        return COOK.User.token;
    }
    logError(fnName, err) {
        if (err && (err.code === undefined || err.code !== 0)) {
            const desc = `Function: ${fnName} | Message: ${err.message} | Details: ${err.detail}`;
            console.error(desc);
            ga('send', 'exception', {
                exDescription: desc,
                exFatal: false,
            });
            return true;
        }
        return false;
    }
}
COOK.Service = new ServiceOld();
function initService() {
    let apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
    if (APP.IsDev) {
        apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
    }
    else if (APP.IsStage) {
        apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
    }
    gapi.client.load('cookservice', 'v1', () => {
        COOK.Service.endpointLoaded();
    }, apiRoot);
    ga('send', {
        hitType: 'timing',
        timingCategory: 'endpoint',
        timingVar: 'init',
        timingValue: window.performance.now(),
    });
}
