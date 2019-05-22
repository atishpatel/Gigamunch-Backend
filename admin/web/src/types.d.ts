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
            signUpDatetimeTimestamp: string;
            activateDatetimeTimestamp: string;
            deactivatedDatetimeTimestamp: string;
        }

        interface ActivityExtended extends Common.Activity {
            dateFull: string;
            status: string;
            addressString: string;
            paidDate: string;
            discountString: string;
        }

        interface LogExtended extends Common.Log {
            color: string;
            icon: string;
            timestampString: string;
            basicPayloadDescriptionHTML: string;
        }
    }
}
