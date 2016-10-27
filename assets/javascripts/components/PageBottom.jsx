/**
 * Created by shudi on 2016/10/21.
 */
import React from 'react';

const propTypes = {
  contents: React.PropTypes.arrayOf(React.PropTypes.string).isRequired,
  styles: React.PropTypes.object,
  height: React.PropTypes.string,
};

function PageBottom({ contents, styles, height }) {
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
        {
          contents.map((content, index) => <p key={`content-${index}`}>{content}</p>)
        }
      </div>
    </div>
  );
}

PageBottom.propTypes = propTypes;

export default PageBottom;