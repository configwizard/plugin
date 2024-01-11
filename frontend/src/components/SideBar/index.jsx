// Sidebar.jsx
import React, { useContext } from 'react';
import { PluginContext } from '../../contexts/plugin.jsx';
import { Link } from 'react-router-dom';
import './styles.css';

const Sidebar = () => {
    const { plugins } = useContext(PluginContext);

    return (
        <div className="sidebar">
            <Link to="/developer">Developer View</Link>
            <Link to="/plugin">Plugin View</Link>
            {plugins.map(p => (
                <div key={p.pluginId} className="sidebar-item">
                    <Link to={`/plugin/${p.pluginId}`}>
                        <div className="plugin-info">
                            <img
                                src={`/${p.pluginId}/icon.svg`}
                                alt={`${p.name} Icon`}
                                className="plugin-icon"
                            />
                            <span>{p.name}</span>
                        </div>
                    </Link>
                </div>
            ))}
        </div>
    );
};

export default Sidebar;
