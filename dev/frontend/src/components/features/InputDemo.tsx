import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

export function InputDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Input</h3>
        <p className="text-sm text-muted-foreground">
          基本的なフォーム入力フィールドを表示します。
        </p>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>基本</CardTitle>
            <CardDescription>標準的なテキスト入力です。</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid w-full max-w-sm items-center gap-1.5">
              <Label htmlFor="email">メールアドレス</Label>
              <Input type="email" id="email" placeholder="Email" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>ファイル</CardTitle>
            <CardDescription>
              ファイルアップロード用の入力です。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid w-full max-w-sm items-center gap-1.5">
              <Label htmlFor="picture">画像</Label>
              <Input id="picture" type="file" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>無効状態</CardTitle>
            <CardDescription>入力が無効化されている状態です。</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid w-full max-w-sm items-center gap-1.5">
              <Label htmlFor="disabled">無効</Label>
              <Input disabled type="email" id="disabled" placeholder="Email" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>ボタン付き</CardTitle>
            <CardDescription>ボタンと組み合わせた入力です。</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex w-full max-w-sm items-center space-x-2">
              <Input type="email" placeholder="Email" />
              <Button type="submit">登録</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
