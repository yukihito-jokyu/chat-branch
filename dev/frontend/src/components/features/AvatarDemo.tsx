import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

export function AvatarDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Avatar</h3>
        <p className="text-sm text-muted-foreground">
          ユーザーを表す画像要素で、画像読み込みエラー時のフォールバック機能を持っています。
        </p>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>基本</CardTitle>
            <CardDescription>画像がある場合のアバターです。</CardDescription>
          </CardHeader>
          <CardContent className="flex items-center gap-4">
            <Avatar>
              <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <Avatar>
              <AvatarImage src="https://github.com/vercel.png" alt="@vercel" />
              <AvatarFallback>VC</AvatarFallback>
            </Avatar>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>フォールバック</CardTitle>
            <CardDescription>
              画像がない場合、フォールバックテキストが表示されます。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center gap-4">
            <Avatar>
              <AvatarImage src="/broken-image.jpg" alt="@broken" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <Avatar>
              <AvatarImage src="" alt="" />
              <AvatarFallback>JD</AvatarFallback>
            </Avatar>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
