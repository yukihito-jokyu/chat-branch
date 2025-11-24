import { useState } from "react";
import { Header } from "@/components/layout/Header";
import { Sidebar } from "@/components/layout/Sidebar";
import { ButtonDemo } from "@/components/features/ButtonDemo";
import { CardDemo } from "@/components/features/CardDemo";
import { InputDemo } from "@/components/features/InputDemo";
import { SelectDemo } from "@/components/features/SelectDemo";
import { DialogDemo } from "@/components/features/DialogDemo";
import { AvatarDemo } from "@/components/features/AvatarDemo";
import { SonnerDemo } from "@/components/features/SonnerDemo";
import { Toaster } from "@/components/ui/sonner";

function App() {
  const [currentView, setCurrentView] = useState("button");

  const renderContent = () => {
    switch (currentView) {
      case "button":
        return <ButtonDemo />;
      case "card":
        return <CardDemo />;
      case "input":
        return <InputDemo />;
      case "select":
        return <SelectDemo />;
      case "dialog":
        return <DialogDemo />;
      case "avatar":
        return <AvatarDemo />;
      case "toast":
        return <SonnerDemo />;
      default:
        return <ButtonDemo />;
    }
  };

  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <Header />
      <div className="container flex-1 items-start md:grid md:grid-cols-[220px_minmax(0,1fr)] md:gap-6 lg:grid-cols-[240px_minmax(0,1fr)] lg:gap-10">
        <Sidebar currentView={currentView} onNavigate={setCurrentView} />
        <main className="relative py-6 lg:gap-10 lg:py-8 xl:grid xl:grid-cols-[1fr_300px]">
          <div className="mx-auto w-full min-w-0">{renderContent()}</div>
        </main>
      </div>
      <Toaster />
    </div>
  );
}

export default App;
