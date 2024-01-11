import {useContext} from "react";
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import './App.css';
import DeveloperView from "./pages/developer/index.jsx";
import PluginView from "./pages/plugin/index.jsx"
import {ContainerProvider} from "./contexts/container.jsx";
import { PluginProvider } from './contexts/plugin.jsx';
import Sidebar from "./components/SideBar/index.jsx";

function App() {

    return (
        <ContainerProvider>
            <PluginProvider>
            <Router>
                <div id="App">
                    <Sidebar></Sidebar>
                    <div className="main-content">
                        <Routes>
                            <Route path="/developer" element={<DeveloperView />} />
                            <Route path="/plugin" element={<PluginView />} />
                            <Route path="/plugin/:id" element={<PluginView />} />
                        </Routes>
                    </div>
                </div>
            </Router>
            </PluginProvider>
        </ContainerProvider>
    )
}

export default App
