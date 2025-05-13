export type Locale = "ja" | "en" | "zh";

export const defaultLocale: Locale = "ja";

export const locales: Record<Locale, string> = {
  ja: "日本語",
  en: "English",
  zh: "中文",
};

export interface TranslationContent {
  meta: {
    title: string;
    description: string;
  };
  nav: {
    languageSelector: string;
  };
  sections: {
    hero: {
      title: string;
      subtitle: string;
    };
    painPoints: {
      title: string;
      points: string[];
    };
    service: {
      title: string;
      description: string;
    };
    howTo: {
      title: string;
      steps: Array<{
        title: string;
        description: string;
      }>;
    };
    notes: {
      title: string;
      points: string[];
    };
    footer: {
      copyright: string;
    };
  };
}

export const translations: Record<Locale, TranslationContent> = {
  ja: {
    meta: {
      title: "UR空き情報モニター - UR賃貸住宅の空き情報をLINEでお知らせ",
      description:
        "UR賃貸住宅の空き物件情報をリアルタイムでLINEに通知します。ご希望の物件が空いたらすぐにお知らせします。",
    },
    nav: {
      languageSelector: "言語を選択",
    },
    sections: {
      hero: {
        title: "UR空き情報モニター",
        subtitle: "UR賃貸住宅の空き情報をLINEでお知らせ",
      },
      painPoints: {
        title: "UR賃貸住宅を探す際の課題",
        points: [
          "人気物件はすぐに埋まってしまう",
          "良い物件を見つけるのが難しい",
          "先着順方式で競争が激しい",
          "毎日希望の物件を確認する必要がある",
          "タイミングを逃すと他の人に予約されてしまう",
        ],
      },
      service: {
        title: "LINEで簡単に空き物件をモニター",
        description:
          "当サービスでは、あなたの希望するUR物件に空きが出た際に、LINEですぐにお知らせします。毎日サイトをチェックする手間を省き、理想の物件を逃さないようサポートします。",
      },
      howTo: {
        title: "使い方",
        steps: [
          {
            title: "1. LINE公式アカウントを友達追加",
            description:
              "まずは当サービスのLINE公式アカウントを友達追加してください。",
          },
          {
            title: "2. 希望物件を登録",
            description:
              "指定の形式でメッセージを送信し、モニターしたい物件を登録します。団地名や間取りなど、希望条件を設定できます。",
          },
          {
            title: "3. 通知を受け取る",
            description:
              "希望条件に合う物件が出た時点で、LINEにて通知が届きます。",
          },
          {
            title: "4. 登録解除方法",
            description:
              "物件のモニターを停止したい場合は、指定のメッセージを送信することで解除できます。",
          },
        ],
      },
      notes: {
        title: "ご利用上の注意点",
        points: [
          "現在、1ユーザーにつき1物件（間取り指定可）のみモニター可能です",
          "空き物件は1回のみ通知されます",
          "希望の物件が見つからない場合は、その物件が貸し出された後に再度登録が必要です",
          "システムの都合上、稀に通知が遅れる場合があります",
        ],
      },
      footer: {
        copyright: "© 2025 UR空き情報モニター",
      },
    },
  },
  en: {
    meta: {
      title:
        "UR Property Monitor - Get LINE notifications for available UR rental housing",
      description:
        "Receive real-time LINE notifications for available UR rental housing. Get notified immediately when your desired property becomes vacant.",
    },
    nav: {
      languageSelector: "Select Language",
    },
    sections: {
      hero: {
        title: "UR Property Monitor",
        subtitle: "Get LINE notifications for available UR rental housing",
      },
      painPoints: {
        title: "Challenges when searching for UR rental housing",
        points: [
          "Popular properties get taken quickly",
          "Difficult to find good properties",
          "First-come-first-served system creates fierce competition",
          "Need to check for desired properties daily",
          "Missing the timing means the property gets reserved by someone else",
        ],
      },
      service: {
        title: "Monitor Available Properties Easily via LINE",
        description:
          "Our service notifies you via LINE immediately when your desired UR property becomes available. Save the hassle of checking the website daily and get support to not miss your ideal property.",
      },
      howTo: {
        title: "How to Use",
        steps: [
          {
            title: "1. Add the LINE Official Account",
            description: "First, add our LINE official account as a friend.",
          },
          {
            title: "2. Register Your Desired Property",
            description:
              "Send a message in the specified format to register the property you want to monitor. You can set conditions such as housing complex name and room layout.",
          },
          {
            title: "3. Receive Notifications",
            description:
              "You will receive a LINE notification when a property matching your desired conditions becomes available.",
          },
          {
            title: "4. How to Unregister",
            description:
              "If you wish to stop monitoring a property, you can do so by sending a specified message.",
          },
        ],
      },
      notes: {
        title: "Notes for Usage",
        points: [
          "Currently, each user can monitor only one property (with room layout specification)",
          "Available properties are notified only once",
          "If you do not find your desired property, you need to register again after that property has been rented out",
          "Due to system limitations, notifications may occasionally be delayed",
        ],
      },
      footer: {
        copyright: "© 2025 UR Property Monitor",
      },
    },
  },
  zh: {
    meta: {
      title: "UR房源监控 - LINE通知UR租赁住房空置信息",
      description:
        "通过LINE实时通知UR租赁住房的空置房源信息。当您心仪的房源一有空置，立即通知您。",
    },
    nav: {
      languageSelector: "选择语言",
    },
    sections: {
      hero: {
        title: "UR房源监控",
        subtitle: "通过LINE通知UR租赁住房空置信息",
      },
      painPoints: {
        title: "寻找UR租赁住房时的挑战",
        points: [
          "热门房源很快就被抢光",
          "很难找到好的房源",
          "先到先得的制度导致竞争激烈",
          "需要每天检查心仪的房源",
          "错过时机意味着房源被他人预定",
        ],
      },
      service: {
        title: "通过LINE轻松监控空置房源",
        description:
          "我们的服务会在您心仪的UR房源空置时通过LINE立即通知您。省去每天查看网站的麻烦，帮助您不错过理想的房源。",
      },
      howTo: {
        title: "使用方法",
        steps: [
          {
            title: "1. 添加LINE官方账号",
            description: "首先，将我们的LINE官方账号添加为好友。",
          },
          {
            title: "2. 注册您心仪的房源",
            description:
              "按指定格式发送消息，注册您想监控的房源。您可以设置诸如住宅小区名称和房间布局等条件。",
          },
          {
            title: "3. 接收通知",
            description: "当符合您期望条件的房源出现时，您将收到LINE通知。",
          },
          {
            title: "4. 如何取消注册",
            description: "如果您希望停止监控房源，可以通过发送指定消息来取消。",
          },
        ],
      },
      notes: {
        title: "使用注意事项",
        points: [
          "目前，每位用户只能监控一个房源（可指定房间布局）",
          "空置房源只通知一次",
          "如果您没有找到心仪的房源，需要在该房源被租出后重新注册",
          "由于系统限制，通知可能偶尔会延迟",
        ],
      },
      footer: {
        copyright: "© 2025 UR房源监控",
      },
    },
  },
};
