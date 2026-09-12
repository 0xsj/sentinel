import { mount } from "svelte";
import App from "./lib/root/App.svelte";

const target = document.getElementById("app");
if (!target) throw new Error("Atelier mount target is missing");

export default mount(App, { target });
