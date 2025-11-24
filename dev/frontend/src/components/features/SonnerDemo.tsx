import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

export function SonnerDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Sonner (Toast)</h3>
        <p className="text-sm text-muted-foreground">
          ユーザーへのフィードバックや通知を表示するスタック可能なトーストコンポーネントです。
        </p>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>基本</CardTitle>
            <CardDescription>
              様々な種類のトースト通知を表示します。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-wrap gap-4">
            <Button
              variant="outline"
              onClick={() => {
                toast("メッセージが送信されました。");
              }}
            >
              シンプル
            </Button>

            <Button
              variant="outline"
              onClick={() => {
                toast("スケジュール完了", {
                  description: "2023年2月10日 金曜日 午後5時57分",
                });
              }}
            >
              タイトル付き
            </Button>

            <Button
              variant="outline"
              onClick={() => {
                toast.error("問題が発生しました", {
                  description: "リクエストに失敗しました。",
                });
              }}
            >
              エラー
            </Button>

            <Button
              variant="outline"
              onClick={() => {
                toast("元に戻しますか？", {
                  action: {
                    label: "元に戻す",
                    onClick: () => console.log("Undo"),
                  },
                });
              }}
            >
              アクション付き
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
