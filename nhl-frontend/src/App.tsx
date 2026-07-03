import "./App.css";
import { Home } from "./home";
import { Teams } from "./teams";
import { Login } from "./login";
import { Search } from "./search";
import { Following } from "./following";
import { Routes, Route } from "react-router-dom";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/following" element={<Following />} />
      <Route path="/home" element={<Home />} />
      <Route path="/search" element={<Search />} />
      <Route path="/teams" element={<Teams />} />
    </Routes>
  );
}

export default App;
