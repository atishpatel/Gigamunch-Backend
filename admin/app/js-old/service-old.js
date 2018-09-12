"use strict";
COOK = COOK || {};
var ServiceOld = (function () {
    function ServiceOld() {
        this.loaded = false;
        this.callQueue = [];
    }
    ServiceOld.prototype.endpointLoaded = function () {
        this.loaded = true;
        this.service = gapi.client.cookservice;
        if (this.callQueue) {
            for (var _i = 0, _a = this.callQueue; _i < _a.length; _i++) {
                var fn = _a[_i];
                fn();
            }
            this.callQueue = [];
        }
        ga('send', {
            hitType: 'timing',
            timingCategory: 'endpoint',
            timingVar: 'load',
            timingValue: window.performance.now(),
        });
    };
    ServiceOld.prototype.schedulePhoneCall = function (phoneNumber, datetime, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.schedulePhoneCall(phoneNumber, datetime, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            datetime: datetime,
            phone_number: phoneNumber,
        };
        this
            .service
            .schedulePhoneCall(request)
            .execute(function (resp) {
            _this.logError('schedulePhoneCall', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.getCook = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getCook(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getCook(request)
            .execute(function (resp) {
            _this.logError('getCook', resp.err);
            callback(resp.cook, resp.err);
        });
    };
    ServiceOld.prototype.updateCook = function (cook, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.updateCook(cook, callback);
            });
            return;
        }
        cook.delivery_price = String(cook.delivery_price);
        var request = {
            gigatoken: this.getToken(),
            cook: cook,
        };
        this
            .service
            .updateCook(request)
            .execute(function (resp) {
            _this.logError('updateCook', resp.err);
            callback(resp.cook, resp.err);
            setTimeout(function () {
                _this.refreshToken();
            }, 1);
        });
    };
    ServiceOld.prototype.getSubMerchant = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getSubMerchant(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubMerchant(request)
            .execute(function (resp) {
            _this.logError('getSubMerchant', resp.err);
            callback(resp.sub_merchant, resp.err);
        });
    };
    ServiceOld.prototype.updateSubMerchant = function (submerchant, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.updateSubMerchant(submerchant, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            sub_merchant: submerchant,
        };
        this
            .service
            .updateSubMerchant(request)
            .execute(function (resp) {
            _this.logError('updateSubMerchant', resp.err);
            callback(resp.cook, resp.err);
            setTimeout(function () {
                _this.refreshToken();
            }, 1);
        });
    };
    ServiceOld.prototype.getMessageToken = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getMessageToken(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            device_id: 'browser',
        };
        this
            .service
            .getMessageToken(request)
            .execute(function (resp) {
            _this.logError('getMessageToken', resp.err);
            callback(resp.token, resp.err);
        });
    };
    ServiceOld.prototype.getInquiry = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getInquiry(id, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            id: id,
        };
        this
            .service
            .getInquiry(request)
            .execute(function (resp) {
            _this.logError('getInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    };
    ServiceOld.prototype.getInquiries = function (startIndex, endIndex, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getInquiries(startIndex, endIndex, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            start_index: startIndex,
            end_index: endIndex,
        };
        this
            .service
            .getInquiries(request)
            .execute(function (resp) {
            _this.logError('getInquiries', resp.err);
            callback(resp.inquiries, resp.err);
        });
    };
    ServiceOld.prototype.acceptInquiry = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.acceptInquiry(id, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            id: id,
        };
        this
            .service
            .acceptInquiry(request)
            .execute(function (resp) {
            _this.logError('acceptInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    };
    ServiceOld.prototype.declineInquiry = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.declineInquiry(id, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            id: id,
        };
        this
            .service
            .declineInquiry(request)
            .execute(function (resp) {
            _this.logError('declineInquiry', resp.err);
            callback(resp.inquiry, resp.err);
        });
    };
    ServiceOld.prototype.getMenus = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getMenus(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getMenus(request)
            .execute(function (resp) {
            _this.logError('getMenus', resp.err);
            callback(resp.menus, resp.err);
        });
    };
    ServiceOld.prototype.saveMenu = function (menu, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.saveMenu(menu, callback);
            });
            return;
        }
        delete menu.items;
        var request = {
            gigatoken: this.getToken(),
            menu: menu,
        };
        this
            .service
            .saveMenu(request)
            .execute(function (resp) {
            _this.logError('saveMenu', resp.err);
            callback(resp.menu, resp.err);
        });
    };
    ServiceOld.prototype.getItem = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getItem(id, callback);
            });
            return;
        }
        var request = {
            id: id,
            gigatoken: this.getToken(),
        };
        this
            .service
            .getItem(request)
            .execute(function (resp) {
            _this.logError('getItem', resp.err);
            callback(resp.item, resp.err);
        });
    };
    ServiceOld.prototype.saveItem = function (item, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.saveItem(item, callback);
            });
            return;
        }
        item.min_servings = Number(item.min_servings);
        item.max_servings = Number(item.max_servings);
        item.cook_price_per_serving = Number(item.cook_price_per_serving);
        var request = {
            gigatoken: this.getToken(),
            item: item,
        };
        this
            .service
            .saveItem(request)
            .execute(function (resp) {
            _this.logError('saveItem', resp.err);
            callback(resp.item, resp.err);
        });
    };
    ServiceOld.prototype.activateItem = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.activateItem(id, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            id: id,
        };
        this
            .service
            .activateItem(request)
            .execute(function (resp) {
            _this.logError('activateItem', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.deactivateItem = function (id, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.deactivateItem(id, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            id: id,
        };
        this
            .service
            .deactivateItem(request)
            .execute(function (resp) {
            _this.logError('deactivateItem', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.getSubLogs = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getSubLogs(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubLogs(request)
            .execute(function (resp) {
            _this.logError('getSubLogs', resp.err);
            callback(resp.sublogs, resp.err);
        });
    };
    ServiceOld.prototype.getSubLogsForDate = function (date, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getSubLogsForDate(date, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
        };
        this
            .service
            .getSubLogsForDate(request)
            .execute(function (resp) {
            _this.logError('getSubLogsForDate', resp.err);
            callback(resp.sublogs, resp.err);
        });
    };
    ServiceOld.prototype.getSubEmails = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getSubEmails(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubEmails(request)
            .execute(function (resp) {
            _this.logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.err);
        });
    };
    ServiceOld.prototype.getSubEmailsAndSubs = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.getSubEmailsAndSubs(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .getSubEmails(request)
            .execute(function (resp) {
            _this.logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.subscribers, resp.err);
        });
    };
    ServiceOld.prototype.skipSubLog = function (date, subEmail, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.skipSubLog(date, subEmail, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
            sub_email: subEmail,
        };
        this
            .service
            .skipSubLog(request)
            .execute(function (resp) {
            _this.logError('skipSubLog', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.CancelSub = function (email, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.CancelSub(email, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            email: email,
        };
        this
            .service
            .CancelSub(request)
            .execute(function (resp) {
            _this.logError('CancelSub', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.discountSubLog = function (date, subEmail, amount, percent, overrideDiscount, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.discountSubLog(date, subEmail, amount, percent, overrideDiscount, callback);
            });
            return;
        }
        var request = {
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
            .execute(function (resp) {
            _this.logError('DiscountSubLog', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.ChangeServingsForDate = function (date, subEmail, servings, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.ChangeServingsForDate(date, subEmail, servings, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            date: date.toISOString(),
            sub_email: subEmail,
            servings: servings,
        };
        this
            .service
            .ChangeServingsForDate(request)
            .execute(function (resp) {
            _this.logError('ChangeServingForDate', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.ChangeServingsPermanently = function (email, servings, vegetarian, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.ChangeServingsPermanently(email, servings, vegetarian, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            email: email,
            servings: servings,
            vegetarian: vegetarian,
        };
        this
            .service
            .ChangeServingsPermanently(request)
            .execute(function (resp) {
            _this.logError('ChangeServingsPermanently', resp.err);
            callback(resp.err);
        });
    };
    ServiceOld.prototype.GetGeneralStats = function (callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.GetGeneralStats(callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .GetGeneralStats(request)
            .execute(function (resp) {
            _this.logError('GetGeneralStats', resp.err);
            callback(resp);
        });
    };
    ServiceOld.prototype.refreshToken = function () {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(this.refreshToken);
            return;
        }
        var request = {
            gigatoken: this.getToken(),
        };
        this
            .service
            .refreshToken(request)
            .execute(function (resp) {
            if (!_this.logError('refreshToken', resp.err)) {
                COOK.User.update(resp.gigatoken);
            }
        });
    };
    ServiceOld.prototype.getToken = function () {
        return COOK.User.token;
    };
    ServiceOld.prototype.logError = function (fnName, err) {
        if (err && (err.code === undefined || err.code !== 0)) {
            var desc = "Function: " + fnName + " | Message: " + err.message + " | Details: " + err.detail;
            console.error(desc);
            ga('send', 'exception', {
                exDescription: desc,
                exFatal: false,
            });
            return true;
        }
        return false;
    };
    return ServiceOld;
}());
function initService() {
}
