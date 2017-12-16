/* exported initService */
/* global COOK gapi ga initService */

declare var COOK: any;
declare var gapi: any;

COOK = COOK || {};

class Service {
  loaded: boolean;
  callQueue: Function[];
  service: any;
  constructor() {
    this.loaded = false;
    this.callQueue = [];
  }

  endpointLoaded() {
    this.loaded = true;
    this.service = gapi.client.cookservice;
    // remove functions from callQueue after calling them.
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

  /*
   * Onboard
   */

  schedulePhoneCall(phoneNumber: string, datetime: string, callback: (err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('schedulePhoneCall', resp.err);
        callback(resp.err);
      });
  }

  finishOnboarding(cook: Cook, submerchant: SubMerchant, callback: (err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('finishOnboarding', resp.err);
        COOK.User.update(resp.gigatoken);
        callback(resp.err);
      });
  }

  /*
   * Cook
   */

  getCook(callback: (cook: Cook, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getCook', resp.err);
        callback(resp.cook, resp.err);
      });
  }

  updateCook(cook: Cook, callback: (cook: Cook, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('updateCook', resp.err);
        callback(resp.cook, resp.err);
        setTimeout(() => {
          this.refreshToken();
        }, 1);
      });
  }

  /*
   * Payout Method
   */

  getSubMerchant(callback: (sub_merchant: SubMerchant, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getSubMerchant', resp.err);
        callback(resp.sub_merchant, resp.err);
      });
  }

  updateSubMerchant(submerchant: SubMerchant, callback: (cook: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('updateSubMerchant', resp.err);
        callback(resp.cook, resp.err);
        setTimeout(() => {
          this.refreshToken();
        }, 1);
      });
  }

  /*
   * Message
   */

  getMessageToken(callback: (token: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getMessageToken', resp.err);
        callback(resp.token, resp.err);
      });
  }

  /*
   * Inquiry
   */

  getInquiry(id: number, callback: (inquiry: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getInquiry', resp.err);
        callback(resp.inquiry, resp.err);
      });
  }

  getInquiries(startIndex: number, endIndex: number, callback: (inquiries: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getInquiries', resp.err);
        // if (window.COOK.isDev) {
        //   callback(this.getFakeInquiries(), resp.err);
        //   return;
        // }
        callback(resp.inquiries, resp.err);
      });
  }

  acceptInquiry(id: string, callback: (inquiry: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('acceptInquiry', resp.err);
        callback(resp.inquiry, resp.err);
      });
  }

  declineInquiry(id: string, callback: (inquiries: any, err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('declineInquiry', resp.err);
        callback(resp.inquiry, resp.err);
      });
  }


  /*
   * Menu
   */

  getMenus(callback: (menus: Menu[], err: ErrorWithCode) => void) {
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
      .execute(
      (resp: Response) => {
        this.logError('getMenus', resp.err);
        callback(resp.menus, resp.err);
      });
  }

  saveMenu(menu: Menu, callback: (cook: any, err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('saveMenu', resp.err);
        callback(resp.menu, resp.err);
      });
  }

  /*
   * Item
   */

  getItem(id: string, callback: (cook: any, err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('getItem', resp.err);
        callback(resp.item, resp.err);
      });
  }

  saveItem(item: Item, callback: (item: Item, err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('saveItem', resp.err);
        callback(resp.item, resp.err);
      });
  }

  activateItem(id: string, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('activateItem', resp.err);
        callback(resp.err);
      });
  }

  deactivateItem(id: string, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('deactivateItem', resp.err);
        callback(resp.err);
      });
  }

  /*
   * Admin
   */

  getSubLogs(callback: (sublogs: SubLogs[], err: ErrorWithCode) => void) {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('getSubLogs', resp.err);
        callback(resp.sublogs, resp.err);
      });
  }

  getSubLogsForDate(date: Date, callback: (sublogs: SubLogs, err: ErrorWithCode) => void) {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('getSubLogsForDate', resp.err);
        callback(resp.sublogs, resp.err);
      });
  }

  getSubEmails(callback: (subEmails: String[], err: ErrorWithCode) => void) {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('getSubEmails', resp.err);
        callback(resp.sub_emails, resp.err);
      });
  }

  skipSubLog(date: Date, subEmail: string, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('skipSubLog', resp.err);
        callback(resp.err);
      });
  }

  CancelSub(email: string, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('CancelSub', resp.err);
        callback(resp.err);
      });
  }

  discountSubLog(date: Date, subEmail: string, amount: number, percent: number, overrideDiscount: boolean, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('DiscountSubLog', resp.err);
        callback(resp.err);
      });
  }

  ChangeServingsForDate(date: Date, subEmail: string, servings: number, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('ChangeServingForDate', resp.err);
        callback(resp.err);
      });
  }

  ChangeServingsPermanently(email: string, servings: number, vegetarian: boolean, callback: (err: ErrorWithCode) => void) {
    // if api is not loaded, add to callQueue
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
      .execute(
      (resp: Response) => {
        this.logError('ChangeServingsPermanently', resp.err);
        callback(resp.err);
      });
  }

  refreshToken() {
    // if api is not loaded, add to _callQueue
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
      .execute(
      (resp: Response) => {
        if (!this.logError('refreshToken', resp.err)) {
          COOK.User.update(resp.gigatoken);
        }
      });
  }

  /*
   * Utils
   */

  getToken() {
    return COOK.User.token;
  }

  logError(fnName: string, err: ErrorWithCode) {
    if (err && (err.code === undefined || err.code !== 0)) {
      const desc = `Function: ${fnName} | Message: ${err.message} | Details: ${err.detail}`;
      console.error(desc);
      ga('send', 'exception', {
        exDescription: desc,
        exFatal: false,
      });
      if (err.code && err.code === 452 && !COOK.isDev) { // code signout
        window.location.href = '/signout';
      }
      return true;
    }
    return false;
  }
}

COOK.Service = new Service();

function initService() {
  let apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
  if (COOK.isDev) {
    //apiRoot = 'http://localhost:8080/_ah/api';
    apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
  } else if (COOK.isStage) {
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

interface ErrorWithCode {
  code: number;
  message: string;
  detail: string;
}

interface Cook {
  delivery_price: string;
}

interface SubMerchant {

}

interface Response {
  gigatoken: string;
  token: string;
  cook: Cook;
  menu: Menu;
  menus: Menu[];
  item: Item;
  items: Item[];
  inquiries: any;
  inquiry: any;
  sublogs: SubLogs[];
  sub_merchant: SubMerchant;
  sub_emails: String[];
  err: ErrorWithCode;
}

interface Menu {
  items: Item[];
}

interface Item {
  min_servings: number;
  max_servings: number;
  cook_price_per_serving: number;
}

interface SubLogs {

}
