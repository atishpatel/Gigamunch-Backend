import Vue from 'vue';

declare global {
    namespace Types {
        interface SubscriberExtended extends Common.Subscriber {
            emails: string[];
            emailsString: string;
            names: string[];
            namesString: string;
            phonenumbers: string[];
            phonenumbersString: string;
            addressString: string;
            addressLink: string;
        }

        interface ActivitiyExtended extends Common.Activity {
            dateFull: string;
        }
    }
}
