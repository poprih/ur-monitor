export type Locale = "ja";

export const defaultLocale: Locale = "ja";

export const locales: Record<Locale, string> = {
  ja: "日本語",
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
};
