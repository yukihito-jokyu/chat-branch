import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function CardDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Card</h3>
        <p className="text-sm text-muted-foreground">
          ヘッダー、コンテンツ、フッターを持つコンテナを表示します。
        </p>
      </div>
      <div className="flex justify-center">
        <Card className="w-[350px]">
          <CardHeader>
            <CardTitle>プロジェクト作成</CardTitle>
            <CardDescription>
              新しいプロジェクトをワンクリックでデプロイします。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="name">名前</Label>
                  <Input id="name" placeholder="プロジェクト名" />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="framework">フレームワーク</Label>
                  <Select>
                    <SelectTrigger id="framework">
                      <SelectValue placeholder="選択してください" />
                    </SelectTrigger>
                    <SelectContent position="popper">
                      <SelectItem value="next">Next.js</SelectItem>
                      <SelectItem value="sveltekit">SvelteKit</SelectItem>
                      <SelectItem value="astro">Astro</SelectItem>
                      <SelectItem value="nuxt">Nuxt.js</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </form>
          </CardContent>
          <CardFooter className="flex justify-between">
            <Button variant="outline">キャンセル</Button>
            <Button variant="outline">デプロイ</Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
