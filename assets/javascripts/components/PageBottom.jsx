/* eslint react/no-array-index-key: ["off"] */
import PropTypes from 'prop-types';
import React from 'react';

export default function PageBottom({ contents, styles, height }) {
  const bottomStyle = {
    position: 'fixed',
    bottom: '0px',
    borderTop: '1px solid #E5E5E5',
    padding: '5px 0 5px 0',
    width: '100%',
    zIndex: 10,
  };

  return (
    <div>
      <div style={{ height }} />
      <div style={{ ...bottomStyle, ...styles }}>
        {contents.map((content, index) => <p key={`content-${index}`}>{content}</p>)}
      </div>
    </div>
  );
}

PageBottom.propTypes = {
  contents: PropTypes.arrayOf(PropTypes.string).isRequired,
  styles: PropTypes.object,
  height: PropTypes.string,
};

PageBottom.defaultProps = {
  styles: {},
  height: '55px',
};
