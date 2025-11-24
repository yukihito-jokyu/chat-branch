import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

export function SelectDemo() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Select</h3>
        <p className="text-sm text-muted-foreground">
          ユーザーがリストから値を選択するためのドロップダウンメニューを表示します。
        </p>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>基本</CardTitle>
            <CardDescription>シンプルなセレクトボックスです。</CardDescription>
          </CardHeader>
          <CardContent>
            <Select>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="テーマを選択" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="light">ライト</SelectItem>
                <SelectItem value="dark">ダーク</SelectItem>
                <SelectItem value="system">システム</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>グループ化</CardTitle>
            <CardDescription>
              項目をグループ化したセレクトボックスです。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Select>
              <SelectTrigger className="w-[280px]">
                <SelectValue placeholder="タイムゾーンを選択" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>北米</SelectLabel>
                  <SelectItem value="est">東部標準時 (EST)</SelectItem>
                  <SelectItem value="cst">中部標準時 (CST)</SelectItem>
                  <SelectItem value="mst">山岳部標準時 (MST)</SelectItem>
                  <SelectItem value="pst">太平洋標準時 (PST)</SelectItem>
                  <SelectItem value="akst">アラスカ標準時 (AKST)</SelectItem>
                  <SelectItem value="hst">ハワイ標準時 (HST)</SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel>ヨーロッパ</SelectLabel>
                  <SelectItem value="gmt">グリニッジ標準時 (GMT)</SelectItem>
                  <SelectItem value="cet">中央ヨーロッパ時間 (CET)</SelectItem>
                  <SelectItem value="eet">東ヨーロッパ時間 (EET)</SelectItem>
                  <SelectItem value="west">
                    西ヨーロッパ夏時間 (WEST)
                  </SelectItem>
                  <SelectItem value="cat">中央アフリカ時間 (CAT)</SelectItem>
                  <SelectItem value="eat">東アフリカ時間 (EAT)</SelectItem>
                </SelectGroup>
                <SelectGroup>
                  <SelectLabel>アジア</SelectLabel>
                  <SelectItem value="msk">モスクワ時間 (MSK)</SelectItem>
                  <SelectItem value="ist">インド標準時 (IST)</SelectItem>
                  <SelectItem value="cst_china">中国標準時 (CST)</SelectItem>
                  <SelectItem value="jst">日本標準時 (JST)</SelectItem>
                  <SelectItem value="kst">韓国標準時 (KST)</SelectItem>
                  <SelectItem value="ist_indonesia">
                    インドネシア中部標準時 (WITA)
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
