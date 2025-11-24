import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function ButtonDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Button</h3>
        <p className="text-sm text-muted-foreground">
          クリック可能な要素を表示し、アクションをトリガーします。
        </p>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-1">
        <Card>
          <CardHeader>
            <CardTitle>バリエーション</CardTitle>
            <CardDescription>
              ボタンのスタイルバリエーションです。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-wrap gap-4">
            <Button variant="default">Default</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="destructive">Destructive</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="ghost">Ghost</Button>
            <Button variant="link">Link</Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>サイズ</CardTitle>
            <CardDescription>
              ボタンのサイズバリエーションです。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <Button size="sm">Small</Button>
            <Button size="default">Default</Button>
            <Button size="lg">Large</Button>
            <Button size="icon">Icon</Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
