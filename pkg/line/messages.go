package line

import (
	"fmt"
)

// MessageTemplates contains all the LINE message templates
var MessageTemplates = struct {
	WelcomeMessage           string
	SubscriptionSuccess      string
	SubscriptionError        string
	UnsubscribeSuccess       string
	UnsubscribeError         string
	InvalidUnitName          string
	DatabaseError           string
	SubscriptionLimitReached string
	SpecifiedRoomTypes       string
	CurrentSubscriptions      string
	InvalidFormat            string
}{
	WelcomeMessage: `Thank you for following us! 

To subscribe to UR property notifications, please send me the exact name of the property you're interested in. I will notify you when vacancies become available. You can subscribe to one property at a time.

You can also specify room types by adding them after the property name with a colon. For example: "恵比寿ビュータワー:3LDK&4LDK" will only notify you about 3LDK and 4LDK units.

To unsubscribe from a property, send "-" followed by the property name. For example: "-恵比寿ビュータワー"

For example, if you want to subscribe to "恵比寿ビュータワー", just send me "恵比寿ビュータワー".

ご利用ありがとうございます！

UR物件の空室通知を受け取るには、ご希望の物件の正確な名称を送信してください。空室が発生した際にお知らせいたします。一度に1物件のみの通知が可能です。

間取りを指定する場合は、物件名の後にコロンと間取りを追加してください。例：「恵比寿ビュータワー:3LDK&4LDK」と送信すると、3LDKと4LDKの空室のみ通知されます。

通知を解除する場合は、「-」の後に物件名を送信してください。例：「-恵比寿ビュータワー」

例えば、「恵比寿ビュータワー」の通知を受け取りたい場合は、「恵比寿ビュータワー」と送信してください。`,
	SubscriptionSuccess: `You have successfully subscribed to UR property %s. You will receive notifications when vacancies become available.

UR %sへの登録が完了しました。空室が発生した際にお知らせいたします。`,
	SubscriptionError: `Failed to subscribe to UR property %s. Please try again later.

UR %sへの登録に失敗しました。しばらくしてから再度お試しください。`,
	UnsubscribeSuccess: `You have successfully unsubscribed from UR property %s. You will no longer receive notifications for this property.

UR %sの通知登録を解除しました。これ以降、この物件の空室通知は送信されません。`,
	UnsubscribeError: `Failed to unsubscribe from UR property %s. Please try again later.

UR %sの通知登録解除に失敗しました。しばらくしてから再度お試しください。`,
	InvalidUnitName: `Invalid property name. Please check the property name and try again.

物件名が正しくありません。正確な物件名を確認の上、再度送信してください。`,
	DatabaseError: `An error occurred while processing your request. Please try again later.

処理中にエラーが発生しました。しばらくしてから再度お試しください。`,
	SubscriptionLimitReached: `Currently, each user can only subscribe to notifications for one property at a time.

現在、お一人様一つの物件のみ空室通知を登録できます。`,
	SpecifiedRoomTypes: `Specified room types: %s

指定された間取り: %s`,
	CurrentSubscriptions: `Current Subscribed Properties:

現在の登録物件:`,
	InvalidFormat: `Please enter in the correct format.
Example: Property Name or Property Name:3LDK&4LDK

正しい形式で入力してください。
例：マンション名 または マンション名:3LDK&4LDK`,
}

// FormatBilingualMessage formats a bilingual message template with the given arguments.
// It automatically duplicates the arguments for both language parts of the message.
func FormatBilingualMessage(template string, args ...interface{}) string {
	// Create a new slice with duplicated arguments
	duplicatedArgs := make([]interface{}, 0, len(args)*2)
	for _, arg := range args {
		duplicatedArgs = append(duplicatedArgs, arg, arg)
	}
	return fmt.Sprintf(template, duplicatedArgs...)
} 
