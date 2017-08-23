"use strict";
/* exported initService */
/* global COOK gapi ga initService */
COOK = COOK || {};
var Service = (function () {
    function Service() {
        this.loaded = false;
        this.callQueue = [];
    }
    Service.prototype.endpointLoaded = function () {
        var _this = this;
        this.loaded = true;
        this.service = gapi.client.cookservice;
        // remove functions from callQueue after calling them.
        if (this.callQueue) {
            for (var _i = 0, _a = this.callQueue; _i < _a.length; _i++) {
                var fn = _a[_i];
                fn();
            }
            this.callQueue = [];
        }
        setTimeout(function () {
            _this.refreshToken();
        }, 3000);
        ga('send', {
            hitType: 'timing',
            timingCategory: 'endpoint',
            timingVar: 'load',
            timingValue: window.performance.now(),
        });
    };
    /*
     * Onboard
     */
    Service.prototype.schedulePhoneCall = function (phoneNumber, datetime, callback) {
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
    Service.prototype.finishOnboarding = function (cook, submerchant, callback) {
        var _this = this;
        if (!this.loaded) {
            this.callQueue.push(function () {
                _this.finishOnboarding(cook, submerchant, callback);
            });
            return;
        }
        var request = {
            gigatoken: this.getToken(),
            cook: cook,
            sub_merchant: submerchant,
        };
        this
            .service
            .finishOnboarding(request)
            .execute(function (resp) {
            _this.logError('finishOnboarding', resp.err);
            COOK.User.update(resp.gigatoken);
            callback(resp.err);
        });
    };
    /*
     * Cook
     */
    Service.prototype.getCook = function (callback) {
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
    Service.prototype.updateCook = function (cook, callback) {
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
    /*
     * Payout Method
     */
    Service.prototype.getSubMerchant = function (callback) {
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
    Service.prototype.updateSubMerchant = function (submerchant, callback) {
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
    /*
     * Message
     */
    Service.prototype.getMessageToken = function (callback) {
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
    /*
     * Inquiry
     */
    Service.prototype.getInquiry = function (id, callback) {
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
    Service.prototype.getInquiries = function (startIndex, endIndex, callback) {
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
            // if (window.COOK.isDev) {
            //   callback(this.getFakeInquiries(), resp.err);
            //   return;
            // }
            callback(resp.inquiries, resp.err);
        });
    };
    Service.prototype.acceptInquiry = function (id, callback) {
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
    Service.prototype.declineInquiry = function (id, callback) {
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
    /*
     * Menu
     */
    Service.prototype.getMenus = function (callback) {
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
    Service.prototype.saveMenu = function (menu, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    /*
     * Item
     */
    Service.prototype.getItem = function (id, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    Service.prototype.saveItem = function (item, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    Service.prototype.activateItem = function (id, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    Service.prototype.deactivateItem = function (id, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    /*
     * Admin
     */
    Service.prototype.getSubLogs = function (callback) {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    Service.prototype.getSubLogsForDate = function (date, callback) {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    Service.prototype.getSubEmails = function (callback) {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    Service.prototype.skipSubLog = function (date, subEmail, callback) {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    Service.prototype.CancelSub = function (email, callback) {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    Service.prototype.discountSubLog = function (date, subEmail, amount, percent, overrideDiscount, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    Service.prototype.ChangeServingsForDate = function (date, subEmail, servings, callback) {
        var _this = this;
        // if api is not loaded, add to callQueue
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
    Service.prototype.refreshToken = function () {
        var _this = this;
        // if api is not loaded, add to _callQueue
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
    /*
     * Utils
     */
    Service.prototype.getToken = function () {
        return COOK.User.token;
    };
    Service.prototype.logError = function (fnName, err) {
        if (err && (err.code === undefined || err.code !== 0)) {
            var desc = "Function: " + fnName + " | Message: " + err.message + " | Details: " + err.detail;
            console.error(desc);
            ga('send', 'exception', {
                exDescription: desc,
                exFatal: false,
            });
            if (err.code && err.code === 452) {
                window.location.href = '/signout';
            }
            return true;
        }
        return false;
    };
    return Service;
}());
COOK.Service = new Service();
function initService() {
    var apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
    if (COOK.isDev) {
        //apiRoot = 'http://localhost:8080/_ah/api';
        apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
    }
    else if (COOK.isStage) {
        apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
    }
    gapi.client.load('cookservice', 'v1', function () {
        COOK.Service.endpointLoaded();
    }, apiRoot);
    ga('send', {
        hitType: 'timing',
        timingCategory: 'endpoint',
        timingVar: 'init',
        timingValue: window.performance.now(),
    });
}
