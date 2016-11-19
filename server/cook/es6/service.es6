/* exported initService */
/* global COOK gapi ga initService */

window.COOK = window.COOK || {};

class Service {
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

  schedulePhoneCall(datetime, callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.schedulePhoneCall(datetime, callback);
      });
      return;
    }
    const request = {
      gigatoken: this.getToken(),
      datetime,
    };
    this
      .service
      .schedulePhoneCall(request)
      .execute(
        (resp) => {
          this.logError('schedulePhoneCall', resp.err);
          callback(resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('finishOnboarding', resp.err);
          COOK.User.update(resp.gigatoken);
          callback(resp.err);
        }
      );
  }

  /*
   * Cook
   */

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
      .execute(
        (resp) => {
          this.logError('getCook', resp.err);
          callback(resp.cook, resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('updateCook', resp.err);
          callback(resp.cook, resp.err);
          setTimeout(() => {
            this.refreshToken();
          }, 1);
        }
      );
  }

  /*
   * Payout Method
   */

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
      .execute(
        (resp) => {
          this.logError('getSubMerchant', resp.err);
          callback(resp.sub_merchant, resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('updateSubMerchant', resp.err);
          callback(resp.cook, resp.err);
          setTimeout(() => {
            this.refreshToken();
          }, 1);
        }
      );
  }

  /*
   * Message
   */

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
      .execute(
        (resp) => {
          this.logError('getMessageToken', resp.err);
          callback(resp.token, resp.err);
        }
      );
  }

  /*
   * Inquiry
   */

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
      .execute(
        (resp) => {
          this.logError('getInquiry', resp.err);
          callback(resp.inquiry, resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('getInquiries', resp.err);
          // if (window.COOK.isDev) {
          //   callback(this.getFakeInquiries(), resp.err);
          //   return;
          // }
          callback(resp.inquiries, resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('acceptInquiry', resp.err);
          callback(resp.inquiry, resp.err);
        }
      );
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
      .execute(
        (resp) => {
          this.logError('declineInquiry', resp.err);
          callback(resp.inquiry, resp.err);
        }
      );
  }


  /*
   * Menu
   */

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
      .execute(
        (resp) => {
          this.logError('getMenus', resp.err);
          callback(resp.menus, resp.err);
        }
      );
  }

  saveMenu(menu, callback) {
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
        (resp) => {
          this.logError('saveMenu', resp.err);
          callback(resp.menu, resp.err);
        }
      );
  }

  /*
   * Item
   */

  getItem(id, callback) {
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
        (resp) => {
          this.logError('getItem', resp.err);
          callback(resp.item, resp.err);
        }
      );
  }

  saveItem(item, callback) {
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
        (resp) => {
          this.logError('saveItem', resp.err);
          callback(resp.item, resp.err);
        }
      );
  }

  activateItem(id, callback) {
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
        (resp) => {
          this.logError('activateItem', resp.err);
          callback(resp.err);
        }
      );
  }

  deactivateItem(id, callback) {
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
        (resp) => {
          this.logError('deactivateItem', resp.err);
          callback(resp.err);
        }
      );
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
        (resp) => {
          if (!this.logError('refreshToken', resp.err)) {
            COOK.User.update(resp.gigatoken);
          }
        }
      );
  }

  /*
   * Utils
   */

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
      if (err.code && err.code === 452) { // code signout
        window.location = '/signout';
      }
      return true;
    }
    return false;
  }
}

window.COOK.Service = new Service();

function initService() {
  let apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
  if (COOK.isDev) {
    apiRoot = 'http://localhost:8080/_ah/api';
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
