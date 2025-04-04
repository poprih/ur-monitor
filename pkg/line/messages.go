package line

// MessageTemplates contains all the LINE message templates
var MessageTemplates = struct {
	WelcomeMessage           string
	SubscriptionSuccess      string
	SubscriptionError        string
	UnsubscribeSuccess       string
	UnsubscribeError         string
	InvalidUnitName          string
	DatabaseError           string
}{
	WelcomeMessage: `Thank you for following us! 

To subscribe to UR property notifications, please send me the exact name of the property you're interested in. I will notify you when vacancies become available. You can subscribe to one property at a time.

ご利用ありがとうございます！

UR物件の空室通知を受け取るには、ご希望の物件の正確な名称を送信してください。空室が発生した際にお知らせいたします。一度に1物件のみの通知が可能です。`,
	SubscriptionSuccess: `You have successfully subscribed to UR %s. You will receive notifications when vacancies become available.

UR%sへの登録が完了しました。空室が発生した際にお知らせいたします。`,
	SubscriptionError: `Failed to subscribe to UR %s. Please try again later.

UR%sへの登録に失敗しました。しばらくしてから再度お試しください。`,
	UnsubscribeSuccess: `You have successfully unsubscribed from UR %s. You will no longer receive notifications for this property.

UR%sの通知登録を解除しました。これ以降、この物件の空室通知は送信されません。`,
	UnsubscribeError: `Failed to unsubscribe from UR %s. Please try again later.

UR%sの通知登録解除に失敗しました。しばらくしてから再度お試しください。`,
	InvalidUnitName: `Invalid unit name. Please check the unit name and try again.

物件名が正しくありません。正確な物件名を確認の上、再度送信してください。`,
	DatabaseError: `An error occurred while processing your request. Please try again later.

処理中にエラーが発生しました。しばらくしてから再度お試しください。`,
} 
