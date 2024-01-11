// Sidebar.jsx
import React, { useContext } from 'react';
import { PluginContext } from '../../contexts/plugin.jsx';
import { Link } from 'react-router-dom';

const Sidebar = () => {
    const { pluginIds } = useContext(PluginContext);

    return (
        <div className="sidebar">
            <Link to="/developer">Developer View</Link>
            <Link to="/plugin">Plugin View</Link>
            {pluginIds.map(id => (
                <Link to={`/plugin/${id}`} key={id}>
                    {id} {/* Displaying the ID for now */}
                </Link>
            ))}
        </div>
    );
};

export default Sidebar;
